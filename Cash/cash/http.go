package cash

import (
	"net/http"
	"strings"
)

const (
	defaultPath = "/__cash__/"
)

type HTTPPool struct {
	self     string
	basePath string
	vs       *viewServer
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self, defaultPath, ViewServer(),
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
	_, _ = w.Write([]byte(group.Info()))
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, pool.basePath) {
		http.NotFound(w, r)
		return
	}

	parts := strings.SplitN(path[len(pool.basePath):], "/", 2)
	groupName := parts[0]
	var group *group
	if group = pool.vs.GetGroup(groupName); group == nil {
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
	return pool.vs.NewGroup(logLevel, maxVolume, namespace, getter)
}

func (pool *HTTPPool) GetGroup(namespace string) *group {
	return pool.vs.GetGroup(namespace)
}
