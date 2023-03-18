package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommentViewCommentReport tells a guest they must be logged in to view reported comments.  */
func (c *GuestController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	c.AddMessage(false, "You must be logged in to see reported comments.")
}

/** HandleCommentViewCommentReport tells a member they don't have permission to view reported comments. */
func (c *MemberController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	c.AddMessage(false, "You don't have permission to see reported comments.")
}

/** HandleCommentViewCommentReport will send back the list of reported comments to the global moderator. */
func (c *GlobalModeratorController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	reports, err := serv.DB.GetCommentReportsFor(comm.Domain)
	if err != nil {
		c.AddMessage(false, fmt.Sprintf("Failed to get comment reports: %v", err.Error()))
	} else {
		c.AddWrapped("CommentReports", reports)
	}
}

/** HandleCommentViewCommentReport will send back the list of reported comments to the domain moderator. */
func (c *DomainModeratorController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	vdomain := false
	for _, d := range c.DomainsModerated {
		if comm.Domain == d {
			vdomain = true
			break
		}
	}
	if !vdomain {
		c.AddMessage(false, "You don't have permission to view reports for that domain.")
	} else {
		reports, err := serv.DB.GetCommentReportsFor(comm.Domain)
		if err != nil {
			c.AddMessage(false, fmt.Sprintf("Failed to get comment reports: %v", err.Error()))
		} else {
			c.AddWrapped("CommentReports", reports)
		}
	}
}

/** HandleCommentViewCommentReport will send back the list of reported comments to the admin. */
func (c *AdminController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	reports, err := serv.DB.GetCommentReportsFor(comm.Domain)
	if err != nil {
		c.AddMessage(false, fmt.Sprintf("Failed to get comment reports: %v", err.Error()))
	} else {
		c.AddWrapped("CommentReports", reports)
	}
}

// viewCommentReports is the API endpoint for when a user attempts to view a list of reported comments. It's called when they send a POST request to "/viewCommentReports". It expects a JSON object of type 'communication.ViewCommentReports'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) viewCommentReports(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewCommentReports{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandViewCommentReports(&comm, s)
		}
	}
}
