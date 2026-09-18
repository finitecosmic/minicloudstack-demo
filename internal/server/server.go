package server

import (
	"minicloudstack/internal/server/handlers"
	"minicloudstack/internal/service/objectstore"
	"minicloudstack/internal/state"
	"net/http"
	"sync/atomic"
)

type Server struct {
	router        http.Handler
	bucketHandler *handlers.BucketHandler
	state         state.State
	ready         atomic.Bool
}

func New() *Server {
	// Create the state implementation.
	memory := state.NewMemory()

	//state injection
	objectstoreService := objectstore.New(memory)

	bucketHandler := handlers.NewBucketHandler(objectstoreService)

	// Create the server with its dependencies.
	s := &Server{
		bucketHandler: bucketHandler,
	}

	// Create the HTTP router and give it access to the server.
	s.router = NewRouter(s)

	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
