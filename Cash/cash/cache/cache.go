package cache

import (
	"fmt"
	"strings"
)

type Value interface {
	Len() int
}

type KVStore interface {
	Get(key string) (Value, bool)
	Put(key string, value Value)
	Del(key string) (Value, bool)
	FlushAll()
	Size() int
}

type Evictor interface {
	Update(key string)
	Del(key string)
	Evict() string
	FlushAll()
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

func ClockAlgorithm(maxVolume int) *Cache {
	return &Cache{
		kvs:       defaultKVStore(),
		evictor:   newClock(),
		maxVolume: maxVolume,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
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

func (c *Cache) Del(key string) (Value, bool) {
	c.evictor.Del(key)
	return c.kvs.Del(key)
}

func (c *Cache) FlushAll() {
	c.evictor.FlushAll()
	c.kvs.FlushAll()
}

func (c *Cache) Info() {
	var percent = float64(c.kvs.Size()) / float64(c.maxVolume)
	numBlocks := int(percent * 20)
	sb := strings.Builder{}
	var i = 0
	for ; i < numBlocks; i++ {
		sb.WriteString("█")
	}
	for ; i < 20; i++ {
		sb.WriteString(" ")
	}
	fmt.Printf("Info of cache:\nCapacity: %d bytes\nUsed %d bytes: %.2f%% |%s|\n", c.maxVolume, c.kvs.Size(), percent*100, sb.String())
}
