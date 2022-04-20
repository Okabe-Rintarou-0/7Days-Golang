package cash

import (
	"Cash/cash/singleflight"
	"fmt"
	"net/http"
	"sync"
)

type Getter interface {
	Get(key string) (ByteView, error)
}

type GetterFunc func(key string) (ByteView, error)

func (f GetterFunc) Get(key string) (ByteView, error) {
	return f(key)
}

type group struct {
	namespace   string
	peerPicker  PeerPicker
	localGetter Getter
	cash        *Cash
	loader      *singleflight.Group
}

func (g *group) Get(key string) (ByteView, error) {
	// Empty key is not allowed.
	if len(key) == 0 {
		return ByteView{}, fmt.Errorf("empty key is not allowed")
	}

	// If the cache hits, then return it.
	if value, ok := g.cash.Get(key); ok {
		return value, nil
	}

	// Otherwise, fetch it by loading.
	return g.load(key)
}

func (g *group) Put(key string, value []byte) error {
	// Empty key is not allowed.
	if len(key) == 0 {
		return fmt.Errorf("empty key is not allowed")
	}

	// If successfully put into a peer, then just return
	if err := g.putInPeer(key, ByteView{value}); err == nil {
		fmt.Printf("Put %s -> %s to the peer\n", key, string(value))
		return nil
	}

	// Put into the cache
	g.cash.Put(key, value)
	return nil
}

func (g *group) Del(key string) (ByteView, error) {
	// Empty key is not allowed.
	if len(key) == 0 {
		return ByteView{}, fmt.Errorf("empty key is not allowed")
	}

	// Del the key in cache
	if value, ok := g.cash.Del(key); ok {
		return value, nil
	} else if value, err := g.delInPeer(key); err == nil {
		fmt.Printf("Delete %s -> %s from a peer\n", key, value.String())
		return value, nil
	}
	return ByteView{}, fmt.Errorf("no such key %s", key)
}

func (g *group) DoBatch(batchedRequest *BatchedRequest) BatchedResponse {
	//baseUrl := fmt.Sprintf("http://%s/__cash__/%s", batchedRequest.Address, batchedRequest.Group)
	numRequests := len(batchedRequest.Requests)
	responses := make([]Response, numRequests)
	wg := sync.WaitGroup{}
	wg.Add(numRequests)
	for id, request := range batchedRequest.Requests {
		switch request.Method {
		case http.MethodPut:
			go func(id int) {
				if err := g.Put(request.Key, []byte(request.Value)); err != nil {
					responses[id] = Response{"Error", err.Error()}
				} else {
					responses[id] = Response{"Ok", ""}
				}
				wg.Done()
			}(id)
		case http.MethodGet:
			go func(id int) {
				if value, err := g.Get(request.Key); err == nil {
					responses[id] = Response{"Ok", value.String()}
				} else {
					responses[id] = Response{"Error", err.Error()}
				}
				wg.Done()
			}(id)
		case http.MethodDelete:
			go func(id int) {
				if value, err := g.Del(request.Key); err == nil {
					responses[id] = Response{"Ok", value.String()}
				} else {
					responses[id] = Response{"Error", err.Error()}
				}
				wg.Done()
			}(id)
		}
	}
	wg.Wait()
	return BatchedResponse{responses}
}

func (g *group) Info() Info {
	return g.cash.Info()
}

func (g *group) FlushAll() {
	g.cash.FlushAll()
}

func (g *group) populate(key string, value ByteView) {
	g.cash.Put(key, value.Clone())
}

func (g *group) load(key string) (ByteView, error) {
	// Try to get from peer, if there exists, then return.
	if value, err := g.loader.DoOnce(key, func() (interface{}, error) {
		return g.getFromPeer(key)
	}); err == nil {
		return value.(ByteView), nil
	}

	var value ByteView
	var err error
	// Otherwise, get locally and populate the key to the cache
	if value, err = g.getLocally(key); err != nil {
		return value, err
	}

	// Populate to the cache.
	g.populate(key, value)
	return value, nil
}

func (g *group) getFromPeer(key string) (ByteView, error) {
	if peer := g.peerPicker.PickPeer(key, g.namespace); peer != nil {
		return peer.Get(key)
	} else {
		return ByteView{}, fmt.Errorf("no available peer")
	}
}

func (g *group) delInPeer(key string) (ByteView, error) {
	if peer := g.peerPicker.PickPeer(key, g.namespace); peer != nil {
		return peer.Del(key)
	} else {
		return ByteView{}, fmt.Errorf("no available peer")
	}
}

func (g *group) putInPeer(key string, value ByteView) error {
	if peer := g.peerPicker.PickPeer(key, g.namespace); peer != nil {
		return peer.Put(key, value)
	} else {
		return fmt.Errorf("no available peer")
	}
}

func (g *group) getLocally(key string) (ByteView, error) {
	return g.localGetter.Get(key)
}

func (g *group) Namespace() string {
	return g.namespace
}
