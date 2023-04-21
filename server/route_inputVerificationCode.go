package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

func (c *GuestController) HandleCommandInputVerificationCode(comm *communication.InputVerificationCode, serv *Server) {
	c.AddMessage(false, "You must be logged in to verify your email.")
}

func (c *MemberControllerBase) HandleCommandInputVerificationCode(comm *communication.InputVerificationCode, serv *Server) {
	ok, msg := serv.DB.AttemptVerify(comm, c.User.ID)
	if ok {
		c.User.IsVerified = true
	}
	prof, err := serv.DB.GetCommUser(c.User)
	if err != nil {
		profResponse := communication.ProfileUpdateResponse{}
		profResponse.Email = c.User.Email
		profResponse.LoggedInAs = *prof
		profResponse.LoggedInAs.IsVerified = c.User.IsVerified
		c.AddWrapped("ProfileUpdateResponse", profResponse)
	}
	c.AddMessage(ok, msg)
}

// postInputVerificationCode is the API endpoint for when a user submits a verificationCode. It expects a JSON object of type 'communication.InputVerificationCode'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postInputVerificationCode(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.InputVerificationCode{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your verification code: %s", err.Error()))
		} else {
			cont.HandleCommandInputVerificationCode(&comm, s)
		}
	}
}
