package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandViewFeedback tells a guest they must be logged in to view reports.  */
func (c *GuestController) HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server) {
	c.AddMessage(false, "You must be logged in to view reports.")
}

/** HandleCommandViewFeedback tells a member they don't have permission to view this report. */
func (c *MemberController) HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server) {
	c.AddMessage(false, "You don't have permission to see the feedback report.")
}

/** HandleCommandViewFeedback tells a member they don't have permission to view this report. */
func (c *GlobalModeratorController) HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server) {
	c.AddMessage(false, "You don't have permission to see the feedback report.")
}

/** HandleCommandViewFeedback tells a member they don't have permission to view this report. */
func (c *DomainModeratorController) HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server) {
	c.AddMessage(false, "You don't have permission to see the feedback report.")
}

/** HandleCommandViewFeedback provides an admin with the feedback report. */
func (c *AdminController) HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server) {
	rep, err := serv.DB.GetFeedbackReport(comm)
	if err != nil {
		c.AddMessage(false, fmt.Sprintf("Failed to get feedbacks: %s", err.Error()))
	} else {
		c.AddWrapped("FeedbackReport", rep)
	}
}

// viewFeedback is the API endpoint for when a user attempts to view a feedback report. It's called when they send a POST request to "/viewFeedback". It expects a JSON object of type 'communication.ViewFeedback'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) viewFeedback(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewFeedback{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandViewFeedback(&comm, s)
		}
	}
}
