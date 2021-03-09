package main

import (
	"go-pg/pkg/postgres"
	"net/http"
)

func main() {
	http.HandleFunc("/api/heros", postgres.GetAllHero)
	http.HandleFunc("/admin/newhero", postgres.CreateHero)
	http.ListenAndServe(":8000", nil)
}