package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandNewReport tells a guest they must be logged in to submit a comment report.  */
func (c *GuestController) HandleCommandNewReport(comm *communication.PostCommentReport, serv *Server) {
	c.AddMessage(false, "You must be logged in to submit a comment report.")
}

// HandleCommandNewReport results in the new comment report being added to the database.
func (c *MemberControllerBase) HandleCommandNewReport(comm *communication.PostCommentReport, serv *Server) {
	c.AddMessage(serv.DB.NewCommentReport(comm, c.User.ID))
}

// postCommentReport is the API endpoint for when a user attempts to add a comment report. It's called when they send a POST request to "/postCommentReport". It expects a JSON object of type 'communication.PostCommentReport'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) postCommentReport(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.PostCommentReport{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandNewReport(&comm, s)
		}
	}
}
