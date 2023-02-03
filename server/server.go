package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/database"
	"github.com/gorilla/mux"
)

type Server struct {
	DB         *database.Store
	httpServer http.Server
	users      UserManager
	router     *mux.Router
}

/* New returns a new server with routing applied and a database connection initialized. */
func New() (*Server, error) {
	if config.Vals.IsLoaded != true {
		return nil, errors.New("Configuration object must be initialized to the create server.")
	}
	s := &Server{}
	s.users = NewUserManager()
	s.users.serv = s

	db, err := database.New(true)
	if err != nil {
		return nil, err
	} else {
		s.DB = db
	}
	s.setupRouter()
	s.httpServer = http.Server{
		Addr: config.Vals.Server.Port,
	}
	return s, nil
}

// setupRouter configures the API endpoints and middleware.
func (s *Server) setupRouter() {
	r := mux.NewRouter()
	// setup the middleware
	r.Use(s.ReadsAuth, s.EnsureController)
	// register api endpoint
	r.HandleFunc("/register", responder(s.postRegister))
	s.router = r
}

/* Start causes the Server to start listening on the port defined in the config.go. */
func (s *Server) Start() {
	fmt.Println("Server listening on port ", config.Vals.Server.Port)
	go fmt.Println(http.ListenAndServe(config.Vals.Server.Port, s.router))

	// If user presses a key, we terminate the server.
	var userIn [10]byte
	os.Stdin.Read(userIn[:])
	fmt.Print(userIn)
}

/* Stop stops the server. */
func (s *Server) Stop() error {
	return s.httpServer.Close()
}

// responder wraps an API endpoint so that it calls "Respond" on the associated controller after all other operations to actually write the response
func responder(last func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		last(w, r)
		cont := r.Context().Value(CtxController).(UserControllerInterface)
		if cont == nil { // if the code is structured correctly, this should never occur, as a guest controller should always be attached to a new request
			w.Write(communication.MessageBytes(false, "Failed to find controller for response!!!"))
		} else {
			cont.Respond(w, r)
		}
	}
}
