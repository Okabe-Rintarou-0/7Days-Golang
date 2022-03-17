package cash

import "fmt"

type logger interface {
	Get(key string, value ByteView)
	Put(key string, bytes []byte)
	Del(key string, value ByteView)
	NotFound(key string)
}

type cashLogger struct {
	logLevel int
}

func defaultLogger(logLevel int) *cashLogger {
	return &cashLogger{logLevel}
}

func (logger *cashLogger) Get(key string, value ByteView) {
	if logger.logLevel > 0 {
		fmt.Printf("[Cash]: Get k-v pair (%s -> %s)\n", key, value.String())
	}
}

func (logger *cashLogger) Put(key string, bytes []byte) {
	if logger.logLevel > 0 {
		fmt.Printf("[Cash]: Put k-v pair (%s -> %s)\n", key, string(bytes))
	}
}

func (logger *cashLogger) Del(key string, value ByteView) {
	if logger.logLevel > 0 {
		fmt.Printf("[Cash]: Del k-v pair (%s -> %s)\n", key, value.String())
	}
}

func (logger *cashLogger) NotFound(key string) {
	if logger.logLevel > 0 {
		fmt.Printf("[Cash]: key %s not found\n", key)
	}
}
