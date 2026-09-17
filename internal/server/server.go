package server

import (
	"minicloudstack/internal/server/handlers"
	"minicloudstack/internal/service/objectstore"
	"minicloudstack/internal/state"
	"net/http"
)

type Server struct {
	router        http.Handler
	bucketHandler *handlers.BucketHandler
}

func New() *Server {
	// Create the state implementation.
	memory := state.NewMemory()

	//Inject state in service
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
