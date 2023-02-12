package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
	"github.com/comment-anything/ca-back-end/util"
)

// HandleCommandChangeEmail on a GuestController will cause an error message to be appended to the next response of the controller.
func (c *GuestController) HandleCommandChangeEmail(comm *communication.ChangeEmail, server *Server) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "You must be logged in to change your email."))
}

// HandleCommandChangeEmail on a UserController will attempt to change the user's email in the database and push an appropriate message depending on the result.
func (c *MemberControllerBase) HandleCommandChangeEmail(comm *communication.ChangeEmail, serv *Server) {

	ok, why := util.ValidateEmail(comm.NewEmail)

	if comm.Password != c.User.Password {
		c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "You must supply your password to change your email."))
	} else if !ok {
		c.AddMessage(ok, why)
	} else {
		params := generated.UpdateUserEmailParams{}
		params.ID = c.User.ID
		params.Email = comm.NewEmail
		err := serv.DB.Queries.UpdateUserEmail(context.Background(), params)
		if err != nil {
			c.AddMessage(false, "Failed to change email.")
		} else {
			c.AddMessage(true, "Email updated.")
			c.User.Email = comm.NewEmail
			prof := serv.GetProfile(c.User)
			profResponse := communication.ProfileUpdateResponse{}
			profResponse.Email = comm.NewEmail
			profResponse.LoggedInAs = prof
			c.AddWrapped("ProfileUpdateResponse", profResponse)
		}
	}
}

// postChangeEmail is the API endpoint for when a user attempts to change their email. It expects a JSON object of type 'communication.Email'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postChangeEmail(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ChangeEmail{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your change email request: %s", err.Error()))
		} else {
			cont.HandleCommandChangeEmail(&comm, s)
		}
	}
}
