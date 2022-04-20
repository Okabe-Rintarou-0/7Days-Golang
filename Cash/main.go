package main

import (
	"Cash/cash"
	"Cash/utils"
	"flag"
	"fmt"
)

func naiveGetter() func(key string) (cash.ByteView, error) {
	kvs := map[string]cash.ByteView{}
	return func(key string) (cash.ByteView, error) {
		if value, ok := kvs[key]; ok {
			return value, nil
		} else {
			return value, fmt.Errorf("no such key")
		}
	}
}

func main() {
	volume4G := 1 << 32
	//volume1G := 1 << 30
	//volume20B := 20
	var maxVolume int
	var logLevel int
	var port int
	flag.IntVar(&maxVolume, "v", 2000, "max-volume")
	flag.IntVar(&logLevel, "l", 2, "log-level")
	flag.IntVar(&port, "p", 8000, "port")
	flag.Parse()

	maxVolume = utils.Clamp(maxVolume, 0, volume4G)
	logLevel = utils.Clamp(logLevel, 0, 2)

	peers := []string{"localhost:8000", "localhost:8001", "localhost:8002"}
	peerMap := map[int]string{
		8000: "localhost:8001",
		8001: "localhost:8002",
		8002: "localhost:8003",
	}

	var self string
	var ok bool
	if self, ok = peerMap[port]; !ok {
		self = "localhost:8001"
	}

	pool := cash.NewHTTPPool(self, peers)
	pool.NewGroup(logLevel, maxVolume, "country", cash.GetterFunc(naiveGetter()))
	pool.NewGroup(logLevel, maxVolume, "game", cash.GetterFunc(naiveGetter()))
	pool.Run()
}
