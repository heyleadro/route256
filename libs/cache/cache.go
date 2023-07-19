package cache

import "route256/libs/cache/lrucache"

type LRUCache interface {
	Add(key interface{}, val interface{})
	Get(key interface{}) (interface{}, bool)
	Len() int
	Del(val interface{}) error
}

type Cache struct {
	LRUCache
}

func NewCache(cap int) *Cache {
	return &Cache{
		LRUCache: lrucache.NewLRUCache(cap),
	}
}
