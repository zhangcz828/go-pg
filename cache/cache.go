package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"sync"
)

// key in global cache
const (
	HeroList = "heroList"
)

var cache *lru.Cache

// the length of cache map
var cacheSize int = 2

var mu sync.Mutex

// 单例模式
func GetCache() *lru.Cache {
	mu.Lock()
	defer mu.Unlock()

	if cache == nil {
		cache, _ = lru.New(cacheSize)
	}

	return cache
}
