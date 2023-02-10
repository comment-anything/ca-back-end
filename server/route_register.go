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

// validateRegisterRequest validates whether a register request can succeed or not by verifying the password, username, etc. are acceptable. If it can't succeed, this function will return a string intended for transmission to the end user explaining generally why the request doesn't work.
func validateRegisterRequest(comm *communication.Register, server *Server) (bool, string) {
	if !comm.AgreedToTerms {
		return false, "You must agree to the terms and conditions."
	}
	if comm.Password != comm.RetypePassword {
		return false, "Passwords must match."
	}
	passGood, why := util.ValidatePasswordStrength(comm.Password)
	if !passGood {
		return false, why
	}
	usernameGood, why := util.ValidateUsername(comm.Username)
	if !usernameGood {
		return false, why
	}
	emailGood, why := util.ValidateEmail(comm.Email)
	if !emailGood {
		return false, why
	}
	_, err := server.DB.Queries.GetUserByEmail(context.Background(), comm.Email)
	if err == nil {
		return false, "That email is in use."
	}
	_, err = server.DB.Queries.GetUserByUserName(context.Background(), comm.Username)
	if err == nil {
		return false, "That username is taken."
	}
	return true, ""
}

// HandleCommandRegister on a GuestController will permit registration if the user does not exist.
func (c *GuestController) HandleCommandRegister(comm *communication.Register, server *Server) {

	canRegister, failMsg := validateRegisterRequest(comm, server)
	if !canRegister {
		c.nextResponse = append(c.nextResponse, communication.GetMessage(false, failMsg))
	} else {
		var args generated.CreateUserParams
		args.Username = comm.Username
		args.Email = comm.Email
		args.Password = comm.Password
		user, err := server.DB.Queries.CreateUser(context.Background(), args)
		if err != nil {
			c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "Failed to register."))
		} else {
			c.manager.TransferGuest(c, &user)
			c.nextResponse = append(c.nextResponse, communication.GetMessage(true, "You registered succesfully."))
			var loginResponse communication.LoginResponse
			loginResponse.LoggedInAs.CreatedOn = user.CreatedAt.Unix()
			loginResponse.LoggedInAs.UserId = user.ID
			loginResponse.LoggedInAs.Username = user.Username
			loginResponse.LoggedInAs.ProfileBlurb = user.ProfileBlurb.String
			c.nextResponse = append(c.nextResponse, communication.Wrap("LoginResponse", loginResponse))

		}
	}
}

// HandleCommandRegister on a UserController will fail.
func (c *MemberControllerBase) HandleCommandRegister(comm *communication.Register, server *Server) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "You are already logged in."))
}

// postRegister is the API endpoint for when a user attempts to register a new account. It expects a JSON object of type 'communication.Register'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. PostRegister then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual registration logic.
func (s *Server) postRegister(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Register{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your registration request: %s", err.Error()))
		} else {
			cont.HandleCommandRegister(&comm, s)
		}
	}
}
