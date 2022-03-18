package cache

type mapKVStore struct {
	store map[string]Value
	size  int
}

func defaultKVStore() *mapKVStore {
	return &mapKVStore{
		store: map[string]Value{},
		size:  0,
	}
}

func (dkvs *mapKVStore) Get(key string) (Value, bool) {
	value, ok := dkvs.store[key]
	return value, ok
}

func (dkvs *mapKVStore) Put(key string, value Value) {
	if oldValue, ok := dkvs.store[key]; ok {
		dkvs.size += value.Len() - oldValue.Len()
	} else {
		dkvs.size += value.Len()
	}
	dkvs.store[key] = value
}

func (dkvs *mapKVStore) Del(key string) (Value, bool) {
	var value Value
	var ok bool
	if value, ok = dkvs.Get(key); ok {
		delete(dkvs.store, key)
		dkvs.size -= value.Len()
	}
	return value, ok
}

func (dkvs *mapKVStore) Size() int {
	return dkvs.size
}

func (dkvs *mapKVStore) FlushAll() {
	dkvs.size = 0
	dkvs.store = map[string]Value{}
}
