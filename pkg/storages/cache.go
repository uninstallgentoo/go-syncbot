package storages

import (
	"sync"
	"time"
)

type CachedEntity struct {
	value          interface{}
	CreatedAt      time.Time
	ExpirationTime int64
}

type CacheStorage struct {
	sync.Mutex
	items map[string]*CachedEntity
}

func NewCacheStorage() *CacheStorage {
	items := make(map[string]*CachedEntity)

	cache := CacheStorage{
		items: items,
	}

	return &cache
}

func (cs *CacheStorage) set(key string, value interface{}, duration time.Duration) {
	cs.items[key] = &CachedEntity{
		value:          value,
		ExpirationTime: time.Now().Add(duration).UnixNano(),
		CreatedAt:      time.Now(),
	}
}

func (cs *CacheStorage) get(key string) interface{} {
	if cs.count() > 0 {
		item := cs.items[key]
		if item == nil {
			return nil
		}
		if time.Now().UnixNano() > item.ExpirationTime {
			delete(cs.items, key)
			return nil
		}
		return item.value
	}
	return nil
}

func (cs *CacheStorage) count() int {
	return len(cs.items)
}

func (cs *CacheStorage) Set(key string, value interface{}, duration time.Duration) {
	cs.Lock()
	defer cs.Unlock()
	cs.set(key, value, duration)
}

func (cs *CacheStorage) Get(key string) interface{} {
	cs.Lock()
	defer cs.Unlock()
	return cs.get(key)
}
