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

func (c *Cash) Get(key string) ByteView {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value := c.cache.Get(key); value != nil {
		c.logger.Get(key, value.(ByteView))
		return value.(ByteView)
	}
	c.logger.NotFound(key)
	return Empty()
}

func (c *Cash) Put(key string, bytes []byte) {
	c.lock.Lock()
	c.logger.Put(key, bytes)
	defer c.lock.Unlock()
	c.cache.Put(key, ByteView{bytes})
}

func (c *Cash) Del(key string) ByteView {
	c.lock.Lock()
	defer c.lock.Unlock()
	if value := c.cache.Del(key); value != nil {
		c.logger.Del(key, value.(ByteView))
		return value.(ByteView)
	}
	c.logger.NotFound(key)
	return Empty()
}

func New(maxVolume, logLevel int) *Cash {
	return &Cash{
		cache:  cache.Default(maxVolume),
		lock:   &sync.RWMutex{},
		logger: defaultLogger(logLevel),
	}
}
