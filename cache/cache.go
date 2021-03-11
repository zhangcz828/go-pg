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

var once sync.Once

// 单例模式
func GetCache() *lru.Cache {
	once.Do(func() {
		cache, _ = lru.New(cacheSize)
	})

	return cache
}
