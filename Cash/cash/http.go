package cash

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultPath = "/__cash__/"
)

type HTTPPool struct {
	self     string
	basePath string
	lock     *sync.RWMutex
	groups   map[string]*group
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultPath,
		lock:     &sync.RWMutex{},
		groups:   map[string]*group{},
	}
}

func (pool *HTTPPool) writeBytes(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(bytes)
}

func (pool *HTTPPool) fail(w http.ResponseWriter, errorMsg string) {
	http.Error(w, errorMsg, http.StatusInternalServerError)
}

func (pool *HTTPPool) handleGet(group *group, w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if value, err := group.Get(key); err == nil {
		pool.writeBytes(w, value.Clone())
	} else {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) handlePut(group *group, w http.ResponseWriter, r *http.Request) {
	key, value := r.URL.Query().Get("key"), r.URL.Query().Get("value")
	if err := group.Put(key, []byte(value)); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) handleDel(group *group, w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if value, err := group.Del(key); err == nil {
		pool.writeBytes(w, value.Clone())
	} else {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) info(group *group, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(group.Info()); err != nil {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	w.Header().Set("Access-Control-Allow-Headers", "*") //header的类型
	w.Header().Set("Access-Control-Allow-Method", "*")
	path := r.URL.Path
	if !strings.HasPrefix(path, pool.basePath) {
		http.NotFound(w, r)
		return
	}

	parts := strings.SplitN(path[len(pool.basePath):], "/", 2)

	if parts[0] == "__groups__" {
		pool.writeGroupsInfo(w)
		return
	}

	groupName := parts[0]
	var group *group
	if group = pool.GetGroup(groupName); group == nil {
		pool.fail(w, "No such group")
		return
	}

	if len(parts) == 2 && parts[1] == "info" {
		pool.info(group, w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		pool.handleGet(group, w, r)
	case http.MethodPut:
		pool.handlePut(group, w, r)
	case http.MethodDelete:
		pool.handleDel(group, w, r)
	}
}

func (pool *HTTPPool) NewGroup(logLevel, maxVolume int, namespace string, getter Getter) *group {
	newGroup := &group{
		namespace: namespace,
		getter:    getter,
		cash:      defaultCash(logLevel, maxVolume, namespace),
	}
	pool.lock.Lock()
	defer pool.lock.Unlock()
	pool.groups[namespace] = newGroup
	return newGroup
}

func (pool *HTTPPool) GetGroup(namespace string) *group {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	return pool.groups[namespace]
}

func (pool *HTTPPool) GroupInfo() []string {
	pool.lock.RLock()
	defer pool.lock.RUnlock()
	var namespaces []string
	for namespace := range pool.groups {
		namespaces = append(namespaces, namespace)
	}
	return namespaces
}

func (pool *HTTPPool) writeGroupsInfo(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(pool.GroupInfo()); err != nil {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) Run() {
	log.Fatal(http.ListenAndServe(pool.self, pool))
}
