package server

import (
	"encoding/json"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

func (c *GuestController) HandleCommandAmILoggedIn(comm *communication.AmILoggedIn, s *Server) {
	s.PageManager.MoveGuestToPage(c, comm.Url, s)
	var logoutResponse communication.LogoutResponse
	c.AddWrapped("LogoutResponse", logoutResponse)
}

func (c *MemberController) HandleCommandAmILoggedIn(comm *communication.AmILoggedIn, s *Server) {
	prof, err := serv.DB.GetCommUser(c.User)
	if err != nil {
		c.AddMessage(false, "There was some problem with your account.")
		return
	}
	s.PageManager.MoveMemberToPage(c, comm.Url, s)
	profResponse := communication.LoginResponse{}
	profResponse.Email = c.User.Email
	profResponse.LoggedInAs = *prof
	c.AddWrapped("LoginResponse", profResponse)
}
func (c *DomainModeratorController) HandleCommandAmILoggedIn(comm *communication.AmILoggedIn, s *Server) {
	prof, err := serv.DB.GetCommUser(c.User)
	if err != nil {
		c.AddMessage(false, "There was some problem with your account.")
		return
	}
	s.PageManager.MoveMemberToPage(c, comm.Url, s)
	profResponse := communication.LoginResponse{}
	profResponse.Email = c.User.Email
	profResponse.LoggedInAs = *prof
	c.AddWrapped("LoginResponse", profResponse)
}
func (c *GlobalModeratorController) HandleCommandAmILoggedIn(comm *communication.AmILoggedIn, s *Server) {
	prof, err := serv.DB.GetCommUser(c.User)
	if err != nil {
		c.AddMessage(false, "There was some problem with your account.")
		return
	}
	s.PageManager.MoveMemberToPage(c, comm.Url, s)
	profResponse := communication.LoginResponse{}
	profResponse.Email = c.User.Email
	profResponse.LoggedInAs = *prof
	c.AddWrapped("LoginResponse", profResponse)
}
func (c *AdminController) HandleCommandAmILoggedIn(comm *communication.AmILoggedIn, s *Server) {
	prof, err := serv.DB.GetCommUser(c.User)
	if err != nil {
		c.AddMessage(false, "There was some problem with your account.")
		return
	}
	s.PageManager.MoveMemberToPage(c, comm.Url, s)
	profResponse := communication.LoginResponse{}
	profResponse.Email = c.User.Email
	profResponse.LoggedInAs = *prof
	c.AddWrapped("LoginResponse", profResponse)
}

func (s *Server) postAmILoggedIn(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.AmILoggedIn{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, "I couldn't see if you were logged in!")
		} else {
			cont.HandleCommandAmILoggedIn(&comm, s)
		}
	}
}
