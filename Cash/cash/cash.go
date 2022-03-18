package cash

import (
	"Cash/cash/cache"
	"sync"
)

type Cash struct {
	cache  *cache.Cache
	lock   *sync.RWMutex
	logger logger
}

func (c *Cash) Get(key string) (ByteView, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.cache.Get(key); ok {
		c.logger.Get(key, value)
		return value.(ByteView), true
	}
	return ByteView{}, false
}

func (c *Cash) Put(key string, bytes []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Put(key, ByteView{bytes})
	c.logger.Put(key, bytes)
}

func (c *Cash) Del(key string) (ByteView, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if value, ok := c.cache.Del(key); ok {
		c.logger.Del(key, value)
		return value.(ByteView), true
	}
	return ByteView{}, false
}

func (c *Cash) Info() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.cache.Info()
}

func (c *Cash) FlushAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.FlushAll()
}

func defaultCash(logLevel, maxVolume int, namespace string) *Cash {
	return &Cash{
		cache:  cache.Default(maxVolume),
		lock:   &sync.RWMutex{},
		logger: defaultLogger(namespace, logLevel),
	}
}

func clockAlgorithmCash(logLevel, maxVolume int, namespace string) *Cash {
	return &Cash{
		cache:  cache.ClockAlgorithm(maxVolume),
		lock:   &sync.RWMutex{},
		logger: defaultLogger(namespace, logLevel),
	}
}
