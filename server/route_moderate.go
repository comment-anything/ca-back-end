package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandModerate tells a guest they don't have permission.  */
func (c *GuestController) HandleCommandModerate(comm *communication.Moderate, serv *Server) {
	c.AddMessage(false, "You do not have permission to moderate a comment.")
}

/** HandleCommandModerate tells a member they don't have permission.  */
func (c *MemberController) HandleCommandModerate(comm *communication.Moderate, serv *Server) {
	c.AddMessage(false, "You do not have permission to moderate a comment.")
}

/** HandleCommandModerate adds a moderation record to the database if the comment is one the domain moderator is permitted to moderate.*/
func (c *DomainModeratorController) HandleCommandModerate(comm *communication.Moderate, serv *Server) {
	str, err := serv.DB.Queries.GetCommentDomain(context.Background(), comm.CommentID)
	if err != nil {
		c.AddMessage(false, err.Error())
	} else if !str.Valid {
		c.AddMessage(false, "Comment did not seem to have a domain.")
	} else {
		for _, val := range c.DomainsModerated {
			if val == str.String {
				c.AddMessage(serv.PageManager.ModerateComment(c.User.ID, comm, serv))
				return
			}
		}
		c.AddMessage(false, fmt.Sprintf("You do not have permission to moderate %s", str.String))
	}
}

/** HandleCommandModerate moderates the comment. It adds the moderation action to the database. */
func (c *GlobalModeratorController) HandleCommandModerate(comm *communication.Moderate, serv *Server) {
	c.AddMessage(serv.PageManager.ModerateComment(c.User.ID, comm, serv))
}

/** HandleCommandModerate moderates the comment. */
func (c *AdminController) HandleCommandModerate(comm *communication.Moderate, serv *Server) {
	c.AddMessage(serv.PageManager.ModerateComment(c.User.ID, comm, serv))
}

// postModerate is the API endpoint for when a user attempts to moderate a comment in response to a report. It's called when they send a POST request to "/moderate". It expects a JSON object of type 'communication.Moderate'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) postModerate(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Moderate{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandModerate(&comm, s)
		}
	}
}
