package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandViewUsersReport tells a guest they must be logged in to see reports.  */
func (c *GuestController) HandleCommandViewUsersReport(comm *communication.ViewUsersReport, serv *Server) {
	c.AddMessage(false, "You must be logged in to view reports.")
}

/** HandleCommandViewUsersReport tells a member they don't have permission to view this report. */
func (c *MemberController) HandleCommandViewUsersReport(comm *communication.ViewUsersReport, serv *Server) {
	c.AddMessage(false, "You don't have permission to see the users report.")
}

/** HandleCommandViewUsersReport provides an admin with the users report. */
func (c *AdminController) HandleCommandViewUsersReport(comm *communication.ViewUsersReport, serv *Server) {
	rep := serv.DB.GetUserReportDBPartial(comm)
	rep.LoggedInUserCount = int64(len(serv.users.members))
	c.AddWrapped("AdminUsersReport", *rep)
}

// viewUsersReport is the API endpoint for when a user attempts to view the user's report. It's called when they send a POST request to "/viewUsersReport". It expects a JSON object of type 'communication.ViewUsersReport'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) viewUsersReport(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewUsersReport{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandViewUsersReport(&comm, s)
		}
	}
}
