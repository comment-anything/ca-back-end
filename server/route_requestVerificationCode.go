package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

func (c *GuestController) HandleCommandRequestVerificationCode(comm *communication.RequestVerificationCode, serv *Server) {
	c.AddMessage(false, "You must be logged in to request a verification code.")

}

func (c *MemberControllerBase) HandleCommandRequestVerificationCode(comm *communication.RequestVerificationCode, serv *Server) {
	if c.User.IsVerified {
		c.AddMessage(false, "You are already verified!")
		return
	}
	c.AddMessage(serv.DB.VerifyRequest(comm, c.User.Email))
}

// postRequestVerificationCode is the API endpoint for when a user submits a password reset request. It expects a JSON object of type 'communication.RequestVerificationCode'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postRequestVerificationCode(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.RequestVerificationCode{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your verification code request: %s", err.Error()))
		} else {
			cont.HandleCommandRequestVerificationCode(&comm, s)
		}
	}
}
