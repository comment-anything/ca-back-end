package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// HandleCommandLogin on a MemberController will fail.
func (c *UserControllerBase) HandleCommandGetComments(comm *communication.GetComments, server *Server) {
	c.AddMessage(false, "Get Comments not implemented.")
}

// postRegister is the API endpoint for when a user attempts to login to their account. It expects a JSON object of type 'communication.Login'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. postLogin then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual login logic.
func (s *Server) getComments(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.GetComments{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your getcomments request: %s", err.Error()))
		} else {
			cont.HandleCommandGetComments(&comm, s)
		}
	}
}
