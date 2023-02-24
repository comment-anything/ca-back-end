package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** ToggleFeedbackHidden tells a guest they must be logged in alter feedback.  */
func (c *GuestController) ToggleFeedbackHidden(comm *communication.ToggleFeedbackHidden, serv *Server) {
	c.AddMessage(false, "You must be logged in to toggle feedback.")
}

/** ToggleFeedbackHidden tells a member they don't have permission to view this report. */
func (c *MemberController) ToggleFeedbackHidden(comm *communication.ToggleFeedbackHidden, serv *Server) {
	c.AddMessage(false, "You don't have permission to alter feedbacks.")
}

/** ToggleFeedbackHidden provides an admin with the feedback report. */
func (c *AdminController) ToggleFeedbackHidden(comm *communication.ToggleFeedbackHidden, serv *Server) {
	ok, msg := serv.DB.ToggleFeedback(comm)
	c.AddMessage(ok, msg)
}

// viewFeedback is the API endpoint for when a user attempts to view a feedback report. It's called when they send a POST request to "/viewFeedback". It expects a JSON object of type 'communication.ViewFeedback'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) toggleFeedbackHidden(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ToggleFeedbackHidden{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.ToggleFeedbackHidden(&comm, s)
		}
	}
}
