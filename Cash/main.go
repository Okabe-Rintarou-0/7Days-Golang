package main

import (
	"Cash/cash"
	"Cash/utils"
	"flag"
	"log"
	"net/http"
)

func main() {
	volume4G := 1 << 32
	//volume1G := 1 << 30
	//volume20B := 20

	v := flag.Int("v", 2000, "max-volume")
	l := flag.Int("l", 2, "log-level")
	flag.Parse()

	maxVolume := utils.Clamp(*v, 0, volume4G)
	logLevel := utils.Clamp(*l, 0, 2)

	addr := "localhost:8000"
	pool := cash.NewHTTPPool(addr)
	pool.NewGroup(logLevel, maxVolume, "country", cash.GetterFunc(naiveGetter()))
	log.Fatal(http.ListenAndServe(addr, pool))
}
