package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/util"
	"golang.org/x/crypto/bcrypt"
)

func validateLoginRequest(comm *communication.Login) (bool, string) {
	return util.ValidateUsernameCharacters(comm.Username)
}

// HandleCommandLogin on a GuestController will permit logging in of the user does not exist.
func (c *GuestController) HandleCommandLogin(comm *communication.Login, server *Server) {

	canLogin, failMsg := validateLoginRequest(comm)
	if !canLogin {
		c.AddMessage(false, failMsg)
		return
	}
	user, err := server.DB.Queries.GetUserByUserName(context.Background(), comm.Username)
	if err != nil {
		c.AddMessage(false, "Could not log in with those credentials.")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(comm.Password))
	if err != nil {
		c.AddMessage(false, "Could not log in with those credentials.")
		return
	}
	if user.Banned {
		c.AddMessage(false, "You have been banned from comment anywhere.")
		return
	}
	c.manager.TransferGuest(c, &user)
	tempcon, err := c.manager.AttemptCreateMemberController(user.ID)
	if err == nil {
		fmt.Println("setting page on new member controller")
		server.PageManager.TransferMemberToPage(tempcon, c.GetPage())

	} else {
		fmt.Println("failed to create member controller", user.ID)
	}
	var loginResponse communication.LoginResponse
	prof, err := server.DB.GetCommUser(&user)
	if err != nil {
		c.AddMessage(false, "There was some problem with your profile.")
	}
	loginResponse.LoggedInAs = *prof
	loginResponse.Email = user.Email
	c.AddWrapped("LoginResponse", loginResponse)
	c.AddMessage(true, "Logged in.")

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
