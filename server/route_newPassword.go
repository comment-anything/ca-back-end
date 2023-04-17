package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
	"golang.org/x/crypto/bcrypt"
)

// HandleCommandChangePassword handles when a user submits a password reset code. If the code is valid, it's deleted from the database.
func (c *UserControllerBase) HandleCommandChangePassword(comm *communication.SetNewPass, serv *Server) {
	response := communication.NewPassResponse{}
	response.Success = false
	if comm.Password != comm.RetypePassword {
		response.Text = "Passwords must match."
		c.AddWrapped("NewPassResponse", response)
		return
	}
	user, err := serv.DB.Queries.GetUserByEmail(context.Background(), comm.Email)
	if err != nil {
		response.Text = "Invalid credentials."
		c.AddWrapped("NewPassResponse", response)
		return
	}
	code, err := serv.DB.Queries.GetPWResetCodeEntry(context.Background(), comm.Code)
	if err != nil || code.UserID != user.ID {
		response.Text = "Invalid code."
		c.AddWrapped("NewPassResponse", response)
		return

	}
	params := generated.SetNewUserPasswordParams{}
	params.Email = comm.Email
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(comm.Password), 8)
	if err != nil {
		response.Text = "Invalid password."
		c.AddWrapped("NewPassResponse", response)
		return
	}
	params.Password = string(hashed_pass)
	err = serv.DB.Queries.SetNewUserPassword(context.Background(), params)
	if err != nil {
		response.Text = "Unable to reset password."
		c.AddWrapped("NewPassResponse", response)
		return
	}
	serv.DB.Queries.DeletePreviousPWRestCodesForUser(context.Background(), user.ID)
	response.Text = "Password reset successful."
	response.Success = true
	c.AddWrapped("NewPassResponse", response)
}

// postSetNewPass is the API endpoint for when a user submits a new password. It expects a JSON object of type 'communication.SetNewPass'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postSetNewPass(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.SetNewPass{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your set new password: %s", err.Error()))
		} else {
			cont.HandleCommandChangePassword(&comm, s)
		}
	}
}
