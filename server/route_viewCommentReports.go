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
	c.AddMessage(false, "Not yet implemented.")
}

/** HandleCommentViewCommentReport will send back the list of reported comments to the domain moderator. */
func (c *DomainModeratorController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	c.AddMessage(false, "Not yet implemented.")
}

/** HandleCommentViewCommentReport will send back the list of reported comments to the admin. */
func (c *AdminController) HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server) {
	c.AddMessage(false, "Pseudo Reported Comment")
	cr := &communication.CommentReports {}
	r := communication.CommentReport{}
	r.ID = 1
	r.ReportingUserID = 234
	r.ReportingUsername = "MrReporter"
	r.Comment = communication.Comment {
		UserId     : 1,
		Username   : "BadPerson",
		CommentId  : 5,
		Content    : "I sad bed athiong",
		Factual    : communication.CommentVoteDimension{
			AlreadyVoted: 0,
			Downs: 0,
			Ups: 0,
		},
		Funny      : communication.CommentVoteDimension{
						AlreadyVoted: 0,
						Downs: 0,
						Ups: 0,
					},
		Agree      : communication.CommentVoteDimension{
						AlreadyVoted: 0,
						Downs: 0,
						Ups: 0,
					},
		Hidden     : false,
		Parent     : 0,
		Removed    : false,
		TimePosted : 743829174,
	}
	r.Reason = "This was very bad content therefore i report."
	r.ActionTaken = false
	r.TimeCreated = 432423523
	cr.Reports = append(cr.Reports, r)
	c.AddWrapped("CommentReports", cr)
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
