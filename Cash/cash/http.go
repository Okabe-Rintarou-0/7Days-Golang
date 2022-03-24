package cash

import (
	"Cash/cash/consistentHash"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

const (
	defaultPath        = "/__cash__/"
	defaultNumReplicas = 50
)

type BatchedResponse struct {
	Responses []string `json:"responses"`
}

type BatchedRequestEntry struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Method string `json:"method"`
}
type BatchedRequest struct {
	Address  string                `json:"address"`
	Group    string                `json:"group"`
	Requests []BatchedRequestEntry `json:"requests"`
}
type HTTPPool struct {
	self     string
	basePath string
	lock     *sync.RWMutex
	groups   map[string]*group
	peers    *consistentHash.ConsistentHash
}

func (pool *HTTPPool) pickPeer(key string) string {
	return pool.peers.Get(key)
}

func NewHTTPPool(self string, peers []string) *HTTPPool {
	pool := &HTTPPool{
		self:     self,
		basePath: defaultPath,
		lock:     &sync.RWMutex{},
		groups:   map[string]*group{},
		peers:    consistentHash.Default(defaultNumReplicas, peers),
	}
	return pool
}

func (pool *HTTPPool) writeBytes(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(bytes)
}

func (pool *HTTPPool) fail(w http.ResponseWriter, errorMsg string) {
	http.Error(w, errorMsg, http.StatusInternalServerError)
}

func (pool *HTTPPool) writeGroupsInfo(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(pool.GroupInfo()); err != nil {
		pool.fail(w, err.Error())
	}
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

func (pool *HTTPPool) allowCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	w.Header().Set("Access-Control-Allow-Headers", "*") //header的类型
	w.Header().Set("Access-Control-Allow-Methods", "*")
}

func (pool *HTTPPool) handleBatch(group *group, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	batchedRequest := BatchedRequest{}
	if err := json.Unmarshal(body, &batchedRequest); err == nil {
		w.WriteHeader(http.StatusOK)
		response := group.DoBatch(&batchedRequest)
		pool.json(w, response)
	} else {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) json(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		pool.fail(w, err.Error())
	}
}

func (pool *HTTPPool) info(group *group, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	pool.json(w, group.Info())
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pool.allowCors(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	reqPath := r.URL.Path
	if !strings.HasPrefix(reqPath, pool.basePath) {
		http.NotFound(w, r)
		return
	}

	parts := strings.SplitN(reqPath[len(pool.basePath):], "/", 2)

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

	if len(parts) == 2 {
		switch parts[1] {
		case "info":
			pool.info(group, w)
		case "__batch__":
			pool.handleBatch(group, w, r)
		}
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

func (pool *HTTPPool) NewGroup(logLevel, maxVolume int, namespace string, localGetter Getter) *group {
	newGroup := &group{
		namespace:   namespace,
		peerPicker:  pool,
		localGetter: localGetter,
		cash:        defaultCash(logLevel, maxVolume, namespace),
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

func (pool *HTTPPool) Run() {
	log.Fatal(http.ListenAndServe(pool.self, pool))
}

func (pool *HTTPPool) PickPeer(key string, namespace string) PeerClient {
	if peer := pool.pickPeer(key); len(peer) > 0 && peer != pool.self {
		return &httpPeerClient{
			self:      peer,
			namespace: namespace,
		}
	}
	return nil
}

type httpPeerClient struct {
	namespace string
	self      string
}

func (hpc *httpPeerClient) baseUrl() string {
	return "http://" + path.Join(hpc.self+defaultPath, hpc.namespace)
}

func (hpc *httpPeerClient) Get(key string) (ByteView, error) {
	requestUrl := fmt.Sprintf("%s?key=%s", hpc.baseUrl(), url.QueryEscape(key))

	var res *http.Response
	var err error
	if res, err = http.Get(requestUrl); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return ByteView{body}, nil
		}
	}
	return ByteView{}, err
}

func (hpc *httpPeerClient) Put(key string, value ByteView) error {
	requestUrl := fmt.Sprintf("%s?key=%s&value=%s", hpc.baseUrl(), url.QueryEscape(key), url.QueryEscape(value.String()))

	var err error
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPut, requestUrl, strings.NewReader("")); err != nil {
		return err
	}

	if _, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	return nil
}

func (hpc *httpPeerClient) Del(key string) (ByteView, error) {
	requestUrl := fmt.Sprintf("%s?key=%s", hpc.baseUrl(), key)

	var err error
	var req *http.Request
	if req, err = http.NewRequest(http.MethodDelete, requestUrl, strings.NewReader("")); err != nil {
		return ByteView{}, err
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err == nil {
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			return ByteView{body}, nil
		}
	}
	return ByteView{}, err
}
