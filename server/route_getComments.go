package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandGetComments on a guest controller calls the appropriate functions on a pagemanager and page so that they can maintain which members and which guests are present. */
func (c *GuestController) HandleCommandGetComments(comm *communication.GetComments, serv *Server) {
	if c.Page != nil {
		c.Page.RemoveGuestFromPage(c)
	}
	c.Page = serv.PageManager.LoadPage(comm.Url, serv)
	c.Page.AddGuestToPage(c)
	c.Page.GetComments(c)

}

/** HandleCommandGetComments on a member controller calls the appropriate functions on a pagemanager and page so that they can maintain which members and which guests are present. */
func (c *MemberControllerBase) HandleCommandGetComments(comm *communication.GetComments, serv *Server) {
	if c.Page != nil {
		c.Page.RemoveMemberFromPage(c)
	}
	c.Page = serv.PageManager.LoadPage(comm.Url, serv)
	c.Page.AddMemberToPage(c)
	c.Page.GetComments(c)
}

// getComments is the API endpoint for when a user attempts to login to their account. It expects a JSON object of type 'communication.GetComments'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. getComments then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual login logic.
func (s *Server) getComments(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.GetComments{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your getcomments request: %s", err.Error()))
		} else {
			cont.HandleCommandGetComments(&comm, s)
		}
	}
}
