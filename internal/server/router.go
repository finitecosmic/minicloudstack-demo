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

	mux.HandleFunc("PUT /buckets/{name}", s.bucketHandler.CreateBucket)
	mux.HandleFunc("POST /buckets/", s.bucketHandler.CreateBucket)
	mux.HandleFunc("GET /buckets", s.bucketHandler.ListBuckets)
	mux.HandleFunc("GET /buckets/{name}", s.bucketHandler.GetBucket)
	mux.HandleFunc("DELETE /buckets/{name}", s.bucketHandler.DeleteBucket)
	return mux
}
