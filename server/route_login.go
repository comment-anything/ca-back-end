package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/util"
)

func validateLoginRequest(comm *communication.Login) (bool, string) {
	return util.ValidateUsernameCharacters(comm.Username)
}

// HandleCommandLogin on a GuestController will permit logging in of the user does not exist.
func (c *GuestController) HandleCommandLogin(comm *communication.Login, server *Server) {

	canLogin, failMsg := validateLoginRequest(comm)
	if !canLogin {
		c.nextResponse = append(c.nextResponse, communication.GetMessage(false, failMsg))
	} else {
		user, err := server.DB.Queries.GetUserByUserName(context.Background(), comm.Username)
		if err != nil {
			c.AddMessage(false, "Could not log in with those credentials.")
		} else {
			if user.Password == comm.Password {
				c.manager.TransferGuest(c, &user)
				var loginResponse communication.LoginResponse
				loginResponse.LoggedInAs = server.GetProfile(&user)
				loginResponse.Email = user.Email
				c.AddWrapped("LoginResponse", loginResponse)
				c.AddMessage(true, "Logged in.")
			} else {
				c.AddMessage(false, "Could not log in with those credentials.")
			}
		}
	}
}

// HandleCommandLogin on a MemberController will fail.
func (c *MemberControllerBase) HandleCommandLogin(comm *communication.Login, server *Server) {
	c.AddMessage(false, "You are already logged in.")
}

// postRegister is the API endpoint for when a user attempts to login to their account. It expects a JSON object of type 'communication.Login'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. postLogin then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual login logic.
func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Login{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your login request: %s", err.Error()))
		} else {
			cont.HandleCommandLogin(&comm, s)
		}
	}
}
