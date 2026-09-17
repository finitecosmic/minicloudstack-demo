package main

import (
	"log"
	"minicloudstack/internal/server"
	"net/http"
)

func main() {
	srv := server.New()
	log.Println("CloudStack listening on :8080")

	if err := http.ListenAndServe(":8080", srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
