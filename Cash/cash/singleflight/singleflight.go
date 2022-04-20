package singleflight

import (
	"fmt"
	"sync"
)

type call struct {
	val interface{}
	err error
	wg  sync.WaitGroup
}

type Group struct {
	mtx     sync.Mutex
	callMap map[string]*call
}

func (g *Group) DoOnce(key string, fun func() (interface{}, error)) (interface{}, error) {
	g.mtx.Lock()
	if g.callMap == nil {
		g.callMap = make(map[string]*call)
	}
	if call, ok := g.callMap[key]; ok {
		g.mtx.Unlock()
		call.wg.Wait()
		fmt.Println("wait ended")
		return call.val, call.err
	}
	call := call{}
	call.wg.Add(1)
	g.callMap[key] = &call
	g.mtx.Unlock()

	call.val, call.err = fun()
	call.wg.Done()

	g.mtx.Lock()
	delete(g.callMap, key)
	g.mtx.Unlock()

	return call.val, call.err
}
