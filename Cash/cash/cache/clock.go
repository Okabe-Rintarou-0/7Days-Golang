package cache

import (
	"container/list"
)

type element struct {
	key        string
	accessFlag bool
}

type clock struct {
	cursor *list.Element
	buffer *list.List
	record map[string]*list.Element
}

func (clk *clock) Update(key string) {
	if rec, ok := clk.record[key]; ok {
		rec.Value.(*element).accessFlag = true
	} else {
		clk.buffer.PushBack(&element{key, true})
		clk.record[key] = clk.buffer.Back()
		if clk.cursor == nil {
			clk.cursor = clk.buffer.Back()
		}
	}
}

func (clk *clock) nextCursor() *list.Element {
	if clk.cursor == clk.buffer.Back() {
		return clk.buffer.Front()
	}
	return clk.cursor.Next()
}

func (clk *clock) Evict() string {
	for clk.cursor.Value.(*element).accessFlag {
		clk.cursor.Value.(*element).accessFlag = false
		clk.cursor = clk.nextCursor()
	}
	evicted := clk.cursor
	clk.cursor = clk.nextCursor()
	key := clk.buffer.Remove(evicted).(*element).key
	delete(clk.record, key)
	if clk.buffer.Len() == 0 {
		clk.cursor = nil
	}
	return key
}

func (clk *clock) Del(key string) {
	if rec, ok := clk.record[key]; ok {
		delete(clk.record, key)
		clk.buffer.Remove(rec)
		if clk.buffer.Len() == 0 {
			clk.cursor = nil
		}
	}
}

func (clk *clock) FlushAll() {
	clk.cursor = nil
	clk.buffer = list.New()
	clk.record = map[string]*list.Element{}
}

func newClock() *clock {
	return &clock{
		cursor: nil,
		buffer: list.New(),
		record: map[string]*list.Element{},
	}
}
