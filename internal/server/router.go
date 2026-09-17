package server

import (
	"minicloudstack/internal/server/handlers"
	"net/http"
)

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("PUT /buckets/{name}", s.bucketHandler.CreateBucket)
	mux.HandleFunc("POST /buckets/", s.bucketHandler.CreateBucket)
	mux.HandleFunc("GET /buckets", s.bucketHandler.ListBuckets)
	mux.HandleFunc("GET /buckets/{name}", s.bucketHandler.GetBucket)
	mux.HandleFunc("DELETE /buckets/{name}", s.bucketHandler.DeleteBucket)
	return mux
}
