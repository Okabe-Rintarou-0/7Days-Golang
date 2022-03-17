package main

import (
	"Cash/cash"
	"flag"
)

func clamp(x, min, max int) int {
	if x > max {
		x = max
	} else if x < min {
		x = min
	}
	return x
}

func main() {
	volume4G := 1 << 32
	volume1G := 1 << 30
	volume20B := 20

	v := flag.Int("v", volume1G, "max-volume")
	l := flag.Int("l", 0, "log-level")
	flag.Parse()

	maxVolume := clamp(*v, volume20B, volume4G)
	logLevel := clamp(*l, 0, 1)

	c := cash.New(maxVolume, logLevel)
	c.Put("Beijing", []byte("China"))
	c.Put("Tokyo", []byte("Japan"))

	c.Get("Beijing")
	c.Del("Beijing")
}
