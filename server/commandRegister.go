package server

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// validateRegisterRequest validates whether a register request can succeed or not. If it can't succeed, this function will return a string intended for transmission to the end user explaining generally why the request doesn't work.
func validateRegisterRequest(comm *communication.Register, server *Server) (bool, string) {
	if !comm.AgreedToTerms {
		return false, "You must agree to the terms and conditions."
	}
	if comm.Password != comm.RetypePassword {
		return false, "Passwords must match."
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
		}
	}
}

// HandleCommandRegister on a UserController will fail.
func (c *UserControllerBase) HandleCommandRegister(comm *communication.Register, server *Server) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(false, "You are already logged in."))
}
