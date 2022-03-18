package cache

import (
	"container/list"
	"fmt"
)

func defaultEvictor() *lru {
	return &lru{
		queue:  list.New(),
		record: make(map[string]*list.Element),
	}
}

type lru struct {
	queue  *list.List
	record map[string]*list.Element
}

func (lru *lru) Update(key string) {
	if rec, ok := lru.record[key]; ok {
		lru.queue.MoveToFront(rec)
	} else {
		lru.queue.PushFront(key)
		lru.record[key] = lru.queue.Front()
	}
}

func (lru *lru) Evict() string {
	key := lru.queue.Remove(lru.queue.Back()).(string)
	fmt.Printf("Evict %s\n", key)
	delete(lru.record, key)
	return key
}

func (lru *lru) Del(key string) {
	if rec, ok := lru.record[key]; ok {
		lru.queue.Remove(rec)
		delete(lru.record, key)
	}
}

func (lru *lru) FlushAll() {
	lru.queue = list.New()
	lru.record = map[string]*list.Element{}
}
