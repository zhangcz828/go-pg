package main

import (
	"go-pg/pkg/postgres"
	"net/http"
)

func main() {
	http.HandleFunc("/", postgres.GetAllHero)
	http.ListenAndServe(":8000", nil)
}