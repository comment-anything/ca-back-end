package server

import "github.com/comment-anything/ca-back-end/database"

type Server struct {
	DB database.Store
}

func New() (*Server, error) {
	server := &Server{}
	return server, nil
}
