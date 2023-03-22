package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandViewLogs tells a guest they must be logged in to view reports.  */
func (c *GuestController) HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server) {
	c.AddMessage(false, "You must be an admin to view logs.")
}

/** HandleCommandViewLogs tells a member they don't have permission to view this report. */
func (c *MemberController) HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server) {
	c.AddMessage(false, "You don't have permission to see the feedback report.")
	c.AddMessage(false, "You must be an admin to view logs.")
}

/** HandleCommandViewLogs tells a member they don't have permission to view this report. */
func (c *GlobalModeratorController) HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server) {
	c.AddMessage(false, "You must be an admin to view logs.")
}

/** HandleCommandViewLogs tells a member they don't have permission to view this report. */
func (c *DomainModeratorController) HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server) {
	c.AddMessage(false, "You must be an admin to view logs.")
}

/** HandleCommandViewLogs provides an admin with the logs they requested. */
func (c *AdminController) HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server) {
	logs, err := serv.DB.GetLogs(comm)
	if err != nil {
		c.AddMessage(false, err.Error())
	} else {
		c.AddWrapped("AdminAccessLogs", logs)
	}
}

// viewLogs is the API endpoint for when a user attempts to view a feedback report. It's called when they send a POST request to "/viewLogs". It expects a JSON object of type 'communication.ViewAccessLogs'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) viewLogs(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewAccessLogs{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandViewLogs(&comm, s)
		}
	}
}
