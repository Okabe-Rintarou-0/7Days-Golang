package cash

import (
	"Cash/cash/cache"
	"fmt"
)

type logger interface {
	Get(key string, value cache.Value)
	Put(key string, bytes []byte)
	Del(key string, value cache.Value)
	NotFound(key string)
}

type cashLogger struct {
	namespace string
	logLevel  int
}

func defaultLogger(namespace string, logLevel int) *cashLogger {
	return &cashLogger{namespace, logLevel}
}

func (logger *cashLogger) Get(key string, value cache.Value) {
	if logger.logLevel >= 2 {
		fmt.Printf("[%s|Cash]: Get k-v pair (%s -> %s)\n", logger.namespace, key, value.(ByteView).String())
	}
}

func (logger *cashLogger) Put(key string, bytes []byte) {
	if logger.logLevel >= 2 {
		fmt.Printf("[%s|Cash]: Put k-v pair (%s -> %s)\n", logger.namespace, key, string(bytes))
	}
}

func (logger *cashLogger) Del(key string, value cache.Value) {
	if logger.logLevel >= 2 {
		fmt.Printf("[%s|Cash]: Del k-v pair (%s -> %s)\n", logger.namespace, key, value.(ByteView).String())
	}
}

func (logger *cashLogger) NotFound(key string) {
	if logger.logLevel >= 2 {
		fmt.Printf("[%s|Cash]: key %s not found\n", logger.namespace, key)
	}
}
