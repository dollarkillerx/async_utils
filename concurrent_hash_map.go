package async_utils

import (
	"fmt"
	"sync"
)

type mapItem struct {
	data interface{}
	mu   sync.RWMutex
}

type ConcurrentHashMap struct {
	db map[string]*mapItem
	mu sync.RWMutex
}

func NewConcurrentHashMap() *ConcurrentHashMap {
	return &ConcurrentHashMap{
		db: map[string]*mapItem{},
	}
}

func (c *ConcurrentHashMap) Insert(key string, val interface{}) {
	item, ex := c.db[key]
	if !ex {
		c.db[key] = &mapItem{
			data: val,
		}
		return
	}

	item.mu.Lock()
	defer item.mu.Unlock()

	item.data = val
}

func (c *ConcurrentHashMap) Get(key string) (interface{}, error) {
	item, ex := c.db[key]
	if !ex {
		return nil, fmt.Errorf("not found")
	}

	item.mu.RLock()
	defer item.mu.RUnlock()

	return item.data, nil
}

type IterationFunc = func(key string, val interface{})

func (c *ConcurrentHashMap) Iteration(iterationFunc IterationFunc) {
	for k := range c.db {
		c.db[k].mu.RLock()
		iterationFunc(k, c.db[k].data)
		c.db[k].mu.RUnlock()
	}
}
