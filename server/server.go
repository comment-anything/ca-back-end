package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/database"
	"github.com/gorilla/mux"
)

type Server struct {
	DB         database.Store
	httpServer http.Server
}

/* New returns a new server with routing applied and a database connection initialized. */
func New() (*Server, error) {
	if config.Vals.IsLoaded != true {
		return nil, errors.New("Configuration object must be initialized to the create server.")
	}
	s := &Server{}
	s.setupRouter()
	s.httpServer = http.Server{
		Addr: config.Vals.Server.Port,
	}
	return s, nil
}

// setupRouter configures the API endpoints and middleware.
func (s *Server) setupRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/register", s.postRegister).Methods(http.MethodPost)
	s.httpServer.Handler = r
}

/* Start causes the Server to start listening on the port defined in the config.go. */
func (s *Server) Start() {
	fmt.Println("Server listening on port ", config.Vals.Server.Port)
	s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}
