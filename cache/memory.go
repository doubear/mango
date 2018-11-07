package cache

import (
	"sync"
	"time"

	"github.com/go-mango/mango/contracts"
)

type memoryItem struct {
	value  interface{}
	expire time.Time
}

func (i *memoryItem) Expired() bool {
	return time.Now().Sub(i.expire) > 0
}

type memoryQueue struct {
	values []interface{}
	length int
	mutex  sync.Mutex
}

//MemoryCacher cacher of mango.
type MemoryCacher struct {
	items  map[string]*memoryItem
	queues map[string]*memoryQueue
	ttr    time.Duration
	mutex  sync.Mutex
}

//Get retrieves cached value from cacher.
func (c *MemoryCacher) Get(id string) interface{} {
	if v, ok := c.items[id]; ok {
		if v.Expired() {
			return nil
		}

		return v.value
	}

	return nil
}

//Set stores given value into cacher.
func (c *MemoryCacher) Set(id string, value interface{}, ttl time.Duration) {
	c.items[id] = &memoryItem{value, time.Now().Add(ttl)}
}

//Del deletes stored value by id.
func (c *MemoryCacher) Del(id string) {
	delete(c.items, id)
}

//Push pushs value into queue.
func (c *MemoryCacher) Push(id string, value interface{}) {
	if q, ok := c.queues[id]; ok {
		q.mutex.Lock()
		q.values = append(q.values, value)
		q.length++
		q.mutex.Unlock()
	} else {
		q := &memoryQueue{}
		q.values = append(q.values, value)
		q.length++

		c.queues[id] = q
	}
}

//Pop pops value from queue.
func (c *MemoryCacher) Pop(id string) interface{} {
	if q, ok := c.queues[id]; ok {
		q.mutex.Lock()
		value := q.values[0]
		q.values = q.values[1:]
		q.length--
		q.mutex.Unlock()

		return value
	}

	return nil
}

//Flush clear all cached data.
func (c *MemoryCacher) Flush() {
	c.items = make(map[string]*memoryItem, 0)
}

//GC clear expired data.
func (c *MemoryCacher) GC() {
	for i, v := range c.items {
		if v.Expired() {
			delete(c.items, i)
		}
	}

	time.AfterFunc(c.ttr, func() {
		c.GC()
	})
}

//Memory creates cache provider instance.
func Memory(gc time.Duration) contracts.Cacher {
	return &MemoryCacher{
		make(map[string]*memoryItem, 0),
		make(map[string]*memoryQueue, 0),
		gc,
		sync.Mutex{},
	}
}
