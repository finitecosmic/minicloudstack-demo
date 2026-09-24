package server

import (
	"minicloudstack/internal/server/handlers"
	"net/http"
)

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()

	healthHander := handlers.NewHealthHandler(s.state)
	mux.HandleFunc("GET /health", healthHander.Health)
	mux.HandleFunc("GET /health/live", healthHander.Live)
	mux.HandleFunc("GET /health/ready", healthHander.Ready)

	mux.HandleFunc("PUT /buckets/{name}", s.bucketHandler.Create)
	mux.HandleFunc("GET /buckets", s.bucketHandler.List)
	mux.HandleFunc("GET /buckets/{name}", s.bucketHandler.Get)
	mux.HandleFunc("DELETE /buckets/{name}", s.bucketHandler.Delete)
	return mux
}
