package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// HandleCommandLogin on a GuestController will permit logging in of the user does not exist.
func (c *GuestController) HandleCommandLogout(comm *communication.Logout, server *Server) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "You are not logged in."))
}

// HandleCommandLogin on a MemberController will fail.
func (c *MemberControllerBase) HandleCommandLogout(comm *communication.Logout, server *Server) {
	c.manager.TransferMember(c)
	var logoutResponse communication.LogoutResponse
	c.AddWrapped("LogoutResponse", logoutResponse)
	c.AddMessage(false, "You have logged out.")
}

// postRegister is the API endpoint for when a user attempts to login to their account. It expects a JSON object of type 'communication.Login'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. postLogin then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual login logic.
func (s *Server) putLogout(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Logout{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your logout request: %s", err.Error()))
		} else {
			cont.HandleCommandLogout(&comm, s)
		}
	}
}
