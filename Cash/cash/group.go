package cash

import "fmt"

type Getter interface {
	Get(key string) (ByteView, error)
}

type GetterFunc func(key string) (ByteView, error)

func (f GetterFunc) Get(key string) (ByteView, error) {
	return f(key)
}

type group struct {
	namespace string
	getter    Getter
	cash      *Cash
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
	}
	return ByteView{}, fmt.Errorf("no such key %s", key)
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
	value, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	// Populate to the cache.
	g.populate(key, value)
	return value, nil
}

func (g *group) Namespace() string {
	return g.namespace
}
