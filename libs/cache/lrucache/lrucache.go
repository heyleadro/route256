package lrucache

import (
	"container/list"
	"errors"
	"sync"
)

type LRUCache struct {
	queue    *list.List
	keys     map[interface{}]*Node
	capacity int
	mu       sync.RWMutex
}

type Node struct {
	val  interface{}
	elem *list.Element
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		queue:    list.New(),
		keys:     make(map[interface{}]*Node),
		capacity: cap,
		mu:       sync.RWMutex{},
	}
}

func (c *LRUCache) Add(key interface{}, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.keys[key]
	if ok {
		c.queue.MoveToFront(node.elem)
		return
	}

	if c.queue.Len() == c.capacity {
		c.removeForAdd(key)
	}

	elem := c.queue.PushFront(key)
	c.keys[key] = &Node{
		val:  val,
		elem: elem,
	}
}

func (c *LRUCache) removeForAdd(key interface{}) {
	backElem := c.queue.Back()
	c.queue.Remove(backElem)
	delete(c.keys, backElem.Value)
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.keys[key]
	if !ok {
		return nil, false // not exists
	}

	c.queue.MoveToFront(node.elem)

	return node.val, true // exists
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.queue.Len()
}

func (c *LRUCache) Del(key interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.keys[key]
	if !ok {
		return errors.New("some del")
	}

	c.queue.Remove(elem.elem)
	delete(c.keys, key)

	return nil
}
