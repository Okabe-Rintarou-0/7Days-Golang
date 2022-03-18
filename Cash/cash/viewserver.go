package cash

import "sync"

type viewServer struct {
	lock   *sync.RWMutex
	groups map[string]*group
}

func ViewServer() *viewServer {
	return &viewServer{
		lock:   &sync.RWMutex{},
		groups: map[string]*group{},
	}
}

func (vs *viewServer) NewGroup(logLevel, maxVolume int, namespace string, getter Getter) *group {
	newGroup := &group{
		namespace: namespace,
		getter:    getter,
		cash:      clockAlgorithmCash(logLevel, maxVolume, namespace),
	}
	vs.lock.Lock()
	defer vs.lock.Unlock()
	vs.groups[namespace] = newGroup
	return newGroup
}

func (vs *viewServer) GetGroup(namespace string) *group {
	vs.lock.RLock()
	defer vs.lock.RUnlock()
	return vs.groups[namespace]
}
