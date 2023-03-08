package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandCommentReply tells a guest they must be logged in to post a comment.  */
func (c *GuestController) HandleCommandCommentReply(comm *communication.CommentReply, serv *Server) {
	c.AddMessage(false, "You must be logged in to post a comment.")
}

/** HandleCommandCommentReply on a member controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *MemberController) HandleCommandCommentReply(comm *communication.CommentReply, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't post a comment here.")
	} else {
		ok, msg := c.Page.NewComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

func (c *DomainModeratorController) HandleCommandCommentReply(comm *communication.CommentReply, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't post a comment here.")
	} else {
		ok, msg := c.Page.NewComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

func (c *GlobalModeratorController) HandleCommandCommentReply(comm *communication.CommentReply, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't post a comment here.")
	} else {
		ok, msg := c.Page.NewComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

/** HandleCommandCommentReply on a admin controller calls the appropriate functions on a pagemanager and page to post a new comment. */
func (c *AdminController) HandleCommandCommentReply(comm *communication.CommentReply, serv *Server) {
	if c.Page == nil {
		c.AddMessage(false, "You can't post a comment here.")
	} else {
		ok, msg := c.Page.NewComment(c, comm, serv)
		if !ok {
			c.AddMessage(false, msg)
		}
	}
}

// postCommentReply is the API endpoint for when a user attempts to post a new comment. It's called when they send a POST request to "/newComment". It expects a JSON object of type 'communication.CommentReply'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. postCommentReply then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual comment-post logic.
func (s *Server) postCommentReply(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.CommentReply{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your getcomments request: %s", err.Error()))
		} else {
			cont.HandleCommandCommentReply(&comm, s)
		}
	}
}
