package main

import (
	"Cash/cash"
	"Cash/utils"
	"flag"
	"fmt"
)

func main() {
	volume4G := 1 << 32
	//volume1G := 1 << 30
	//volume20B := 20

	v := flag.Int("v", 10, "max-volume")
	l := flag.Int("l", 2, "log-level")
	flag.Parse()

	maxVolume := utils.Clamp(*v, 0, volume4G)
	logLevel := utils.Clamp(*l, 0, 2)

	vs := cash.ViewServer()
	countries := vs.NewGroup(logLevel, maxVolume, "Country", cash.GetterFunc(naiveGetter()))
	var method, key, value string
	for true {
		n, err := fmt.Scanf("%s", &method)
		if n > 0 && err == nil {
			switch method {
			case "del":
				fmt.Scanf("%s", &key)
				if v, err := countries.Del(key); err == nil {
					fmt.Printf("Del %s -> %s succeeded\n", key, v)
				} else {
					fmt.Printf("Del key %s failed\n", key)
				}
			case "get":
				fmt.Scanf("%s", &key)
				if value, err := countries.Get(key); err == nil {
					fmt.Printf("Get k-v pair (%s -> %s)\n", key, value.String())
				} else {
					fmt.Printf("Get key %s failed\n", key)
				}
			case "put":
				fmt.Scanf("%s %s", &key, &value)
				if err := countries.Put(key, []byte(value)); err == nil {
					fmt.Printf("Put key %s succeeded\n", key)
				} else {
					fmt.Printf("Put key %s failed\n", key)
				}
			case "info":
				countries.Info()
			case "flush":
				countries.FlushAll()
			}
		}
	}
}
