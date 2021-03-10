package main

import (
	"go-pg/cache"
	"go-pg/pkg/postgres"
	"net/http"
)

func main() {
	// 定义缓存节点数量以及地址
	peersAddr := []string{"http://loalhost:8001"}

	cache.SetCacheHttpPool(peersAddr[0], peersAddr)

	http.HandleFunc("/api/heros", postgres.GetAllHeros)
	http.HandleFunc("/admin/newhero", postgres.CreateHero)

	//for _, addr := range peersAddr {
	//	go func() {
	//		cache.SetCacheHttpPool(addr, peersAddr)
	//		http.ListenAndServe(addr, nil)
	//	}()
	//}
	//
	//waitCh := make(chan struct{})
	//
	//<- waitCh

	http.ListenAndServe("localhost:8001", nil)
}