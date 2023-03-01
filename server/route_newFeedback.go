package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
	"github.com/comment-anything/ca-back-end/util"
)

/** HandleCommandNewFeedback tells a guest they must be logged in to submit feedback.  */
func (c *GuestController) HandleCommandNewFeedback(comm *communication.Feedback, serv *Server) {
	c.AddMessage(false, "You must be logged in to submit feedback.")
}

// HandleCommandNewFeedback results in the new feedback being added to the database.
func (c *MemberControllerBase) HandleCommandNewFeedback(comm *communication.Feedback, serv *Server) {
	ok, why := util.ValidateFeedbackType(comm.FeedbackType)
	if !ok {
		c.AddMessage(false, why)
	} else {
		arg := generated.CreateFeedbackParams{}
		arg.Content = comm.Content
		arg.Type = comm.FeedbackType
		arg.UserID = c.User.ID
		err := serv.DB.Queries.CreateFeedback(context.Background(), arg)
		if err != nil {
			c.AddMessage(false, "Failed to submit feedback.")
		} else {
			c.AddMessage(true, "Your feedback has been submitted.")
		}
	}
}

// newFeedback is the API endpoint for when a user attempts to view a feedback report. It's called when they send a POST request to "/newFeedback". It expects a JSON object of type 'communication.Feedback'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) newFeedback(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Feedback{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandNewFeedback(&comm, s)
		}
	}
}
