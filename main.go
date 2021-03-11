package main

import (
	"go-pg/pkg/postgres"
	"net/http"
)

func main() {

	// User rest api
	http.HandleFunc("/api/heros", postgres.GetAllHeros)

	// Admin rest api
	http.HandleFunc("/admin/newhero", postgres.CreateHero)

	http.ListenAndServe("localhost:8000", nil)
}