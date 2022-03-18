package main

import (
	"Cash/cash"
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
