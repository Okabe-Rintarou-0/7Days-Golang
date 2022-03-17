package cache

type Value interface {
	Len() int
}

type KVStore interface {
	Get(key string) Value
	Put(key string, value Value)
	Del(key string) Value
	Size() int
}

type Evictor interface {
	Update(key string)
	Del(key string)
	Evict() string
}

type Cache struct {
	kvs       KVStore
	evictor   Evictor
	maxVolume int
}

func Default(maxVolume int) *Cache {
	return &Cache{
		kvs:       defaultKVStore(),
		evictor:   defaultEvictor(),
		maxVolume: maxVolume,
	}
}

func (c *Cache) Get(key string) Value {
	c.evictor.Update(key)
	return c.kvs.Get(key)
}

func (c *Cache) Put(key string, value Value) {
	c.evictor.Update(key)
	c.kvs.Put(key, value)
	for c.kvs.Size() > c.maxVolume {
		c.kvs.Del(c.evictor.Evict())
	}
}

func (c *Cache) Del(key string) Value {
	c.evictor.Del(key)
	return c.kvs.Del(key)
}
