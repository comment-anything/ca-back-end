package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandVoteComment tells a guest they must be logged in to post a comment.  */
func (c *GuestController) HandleCommandVoteComment(comm *communication.CommentVote, serv *Server) {
	c.AddMessage(false, "You must be logged in to vote on a comment.")
}

/** HandleCommandVoteComment on a member controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *MemberController) HandleCommandVoteComment(comm *communication.CommentVote, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't vote on a comment here.")
	} else {
		ok, msg := c.Page.VoteComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

/** HandleCommandVoteComment on a member controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *DomainModeratorController) HandleCommandVoteComment(comm *communication.CommentVote, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't vote on a comment here.")
	} else {
		ok, msg := c.Page.VoteComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

/** HandleCommandVoteComment on a member controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *GlobalModeratorController) HandleCommandVoteComment(comm *communication.CommentVote, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't vote on a comment here.")
	} else {
		ok, msg := c.Page.VoteComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

/** HandleCommandVoteComment on a member controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *AdminController) HandleCommandVoteComment(comm *communication.CommentVote, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't vote on a comment here.")
	} else {
		ok, msg := c.Page.VoteComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

// voteComment is the API endpoint for when a user attempts to vote on a comment. It's called when they send a POST request to "/voteComment". It expects a JSON object of type 'communication.CommentVote'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. voteComment then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual comment-post logic.
func (s *Server) voteComment(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.CommentVote{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your vote request: %s", err.Error()))
		} else {
			cont.HandleCommandVoteComment(&comm, s)
		}
	}
}
