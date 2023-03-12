package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/database"
	"github.com/gorilla/mux"
)

type Server struct {
	DB          *database.Store
	httpServer  http.Server
	users       UserManager
	PageManager PageManager
	router      *mux.Router
}

/* New returns a new server with routing applied and a database connection initialized. */
func New() (*Server, error) {
	if config.Vals.IsLoaded != true {
		return nil, errors.New("Configuration object must be initialized to the create server.")
	}
	s := &Server{}
	s.users = NewUserManager()
	s.PageManager = NewPageManager()
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
	r.Use(CORS, LogMiddleware, s.ReadsAuth, s.EnsureController)

	// register api endpoint
	r.HandleFunc("/register", responder(s.postRegister))
	r.HandleFunc("/login", responder(s.postLogin))
	r.HandleFunc("/logout", responder(s.putLogout))
	r.HandleFunc("/changeEmail", responder(s.postChangeEmail))
	r.HandleFunc("/changeProfile", responder(s.postChangeProfileBlurb))
	r.HandleFunc("/pwResetReq", responder(s.postPasswordResetRequest))
	r.HandleFunc("/newPassword", responder(s.postSetNewPass))
	r.HandleFunc("/getComments", responder(s.getComments))
	r.HandleFunc("/newComment", responder(s.postCommentReply))
	r.HandleFunc("/voteComment", responder(s.voteComment))
	r.HandleFunc("/voteComment", responder(s.voteComment))
	r.HandleFunc("/viewUsersReport", responder(s.viewUsersReport))
	r.HandleFunc("/viewFeedback", responder(s.viewFeedback))
	r.HandleFunc("/toggleFeedbackHidden", responder(s.toggleFeedbackHidden))
	r.HandleFunc("/newFeedback", responder(s.newFeedback))
	r.HandleFunc("/assignGlobalModerator", responder(s.assignGlobalModerator))
	r.HandleFunc("/assignAdmin", responder(s.assignAdmin))

	s.httpServer.Handler = s.router
	s.router = r
}

/* Start causes the Server to start listening on the port defined in the config.go. Generally, the CLI is not started during tests (so the parameter passed is false). */
func (s *Server) Start(startCli bool) {
	fmt.Println("Server running on port ", config.Vals.Server.Port)
	s.httpServer.Handler = s.router
	s.httpServer.Addr = config.Vals.Server.Port
	if startCli {
		go s.httpServer.ListenAndServe()
		/* The CLI blocks on the main thread. */
		s.CLIBegin()
		/* When the CLI is done, we can close the server. */
		s.httpServer.Close()
	} else {
		/* if startCli is false, server blocks main thread. */
		s.httpServer.ListenAndServe()
	}

}

/* Stop stops the server. */
func (s *Server) Stop() error {
	return s.httpServer.Close()
}

// responder wraps an API endpoint so that it calls "SetCookie" and "Respond" on the associated controller after all other operations to actually write the response and provide the token. It wraps its middleware in the EndLogMiddleware so that it can provide some logging before the request finishes, if desired.
func responder(last func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	responder_func := func(w http.ResponseWriter, r *http.Request) {
		last(w, r)
		cont := r.Context().Value(CtxController).(UserControllerInterface)
		if cont == nil { // if the code is structured correctly, this should never occur, as a guest controller should always be attached to a new request
			w.Write(communication.GetErrMsg(false, "Failed to find controller for response!!!"))
		} else {
			cont.SetCookie(w, r)
			cont.Respond(w, r)
		}
	}
	return EndLogMiddleware(responder_func)
}
