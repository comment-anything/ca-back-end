package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
	"github.com/comment-anything/ca-back-end/util"
)

// HandleCommandChangeProfileBlurb on a GuestController will cause an error message to be appended to the next response of the controller.
func (c *GuestController) HandleCommandChangeProfileBlurb(comm *communication.ChangeProfileBlurb, server *Server) {
	c.AddMessage(false, "You must be logged in to change your profile.")
}

// HandleCommandChangeProfileBlurb on a UserController will attempt to change the user's profile blurb in the database and push an appropriate message depending on the result.
func (c *MemberControllerBase) HandleCommandChangeProfileBlurb(comm *communication.ChangeProfileBlurb, server *Server) {
	ok, why := util.ValidateProfile(comm.NewBlurb)
	if !ok {
		c.AddMessage(false, why)
	} else {
		params := generated.UpdateUserProfileBlurbParams{}
		params.ID = c.User.ID
		params.ProfileBlurb = sql.NullString{String: comm.NewBlurb, Valid: true}
		err := server.DB.Queries.UpdateUserProfileBlurb(context.Background(), params)
		if err != nil {
			c.AddMessage(false, "Unable to change profile blurb.")
		} else {
			c.AddMessage(true, "Profile updated.")
			c.User.ProfileBlurb = params.ProfileBlurb
		}
	}
}

// postChangeProfileBlurb is the API endpoint for when a user attempts to change their profile blurb. It expects a JSON object of type 'communication.ChangeProfileBlurb'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postChangeProfileBlurb(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ChangeProfileBlurb{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your change profile request: %s", err.Error()))
		} else {
			cont.HandleCommandChangeProfileBlurb(&comm, s)
		}
	}
}
