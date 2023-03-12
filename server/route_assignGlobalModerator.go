package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandAssignGlobalModerator tells a guest they dont have permission.  */
func (c *GuestController) HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignGlobalModerator tells a member they don't have permission */
func (c *MemberController) HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignGlobalModerator tells a member they don't have permission */
func (c *GlobalModeratorController) HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignGlobalModerator tells a member they don't have permission. */
func (c *DomainModeratorController) HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server) {
	c.AddMessage(false, "You don't have permission to execute this command.")
}

/** HandleCommandAssignGlobalModerator creates a new global mod assignment in the database. If the user targetted is currently logged in, it will also call methods on UserManager to swap their controller. */
func (c *AdminController) HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server) {
	id, err := serv.DB.GetUserID(comm.User)
	if err != nil {
		c.AddMessage(false, "Couldn't locate that user.")
		return
	}
	err = serv.DB.AssignGlobalModerator(id, c.User.ID, comm.IsDeactivation)
	if err != nil {
		c.AddMessage(false, err.Error())
		return
	}
	if comm.IsDeactivation == false {
		c.AddMessage(serv.users.ChangeMemberControllerToGlobalModController(id))
	} else {
		c.AddMessage(serv.users.RemoveGlobalModPrivileges(id))
	}
}

// assignGlobalModerator is the API endpoint for when a user attempts to assign a global moderator. It's called when they send a POST request to "/assignGlobalModerator". It expects a JSON object of type 'communication.AssignGlobalModerator'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) assignGlobalModerator(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.AssignGlobalModerator{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandAssignGlobalModerator(&comm, s)
		}
	}
}
