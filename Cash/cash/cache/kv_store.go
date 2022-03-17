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

func (dkvs *mapKVStore) Get(key string) Value {
	if value, ok := dkvs.store[key]; ok {
		return value
	}
	return nil
}

func (dkvs *mapKVStore) Put(key string, value Value) {
	if oldValue, ok := dkvs.store[key]; ok {
		dkvs.size += int(value.Len() - oldValue.Len())
	} else {
		dkvs.size += int(value.Len())
	}
	dkvs.store[key] = value
}

func (dkvs *mapKVStore) Del(key string) Value {
	var value Value
	if value = dkvs.Get(key); value != nil {
		delete(dkvs.store, key)
		dkvs.size -= int(value.Len())
	}
	return value
}

func (dkvs *mapKVStore) Size() int {
	return dkvs.size
}
