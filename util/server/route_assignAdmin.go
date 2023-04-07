package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandAssignAdmin tells a guest they dont have permission.  */
func (c *GuestController) HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignAdmin tells a member they don't have permission */
func (c *MemberController) HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignAdmin tells a member they don't have permission */
func (c *GlobalModeratorController) HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignAdmin tells a member they don't have permission. */
func (c *DomainModeratorController) HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignAdmin creates a admin mod assignment in the database. If the user targetted is currently logged in, it will also call methods on UserManager to swap their controller. */
func (c *AdminController) HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server) {
	id, err := serv.DB.GetUserID(comm.User)
	if err != nil {
		c.AddMessage(false, "Couldn't locate that user.")
		return
	}
	err = serv.DB.AssignAdmin(id, c.User.ID)
	if err != nil {
		c.AddMessage(false, err.Error())
		return
	}
	c.AddMessage(serv.users.ChangeMemberControllerToAdminController(id))
}

// assignAdmin is the API endpoint for when a user attempts to assign an admin. It's called when they send a POST request to "/assignAdmin". It expects a JSON object of type 'communication.AssignAdmin'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) assignAdmin(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.AssignAdmin{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandAssignAdmin(&comm, s)
		}
	}
}
