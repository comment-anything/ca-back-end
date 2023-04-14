package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// HandleCommandPasswordResetRequest handles a user's request for a new password by generating a unique code, saving it in the database, and deleting any previous codes for that user.
func (c *UserControllerBase) HandleCommandPasswordResetRequest(comm *communication.PasswordReset, serv *Server) {
	serv.DB.PwResetRequest(comm)
	c.AddMessage(true, "A password reset code has been email to you.")
}

// postPasswordResetRequest is the API endpoint for when a user submits a password reset request. It expects a JSON object of type 'communication.PasswordReset'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.PasswordReset{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your password reset request: %s", err.Error()))
		} else {
			cont.HandleCommandPasswordResetRequest(&comm, s)
		}
	}
}
