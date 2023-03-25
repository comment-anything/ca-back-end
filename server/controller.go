package server

import (
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// UserControllerInterface provides method signatures which other UserController types implement. Controller references are attached to HTTP Request Contexts in the first middleware that a Request passes through. Those controller references are subsequently used by API endpoints to execute access-appropriate code associated with a particular user or guest. At the API endpoints, the Server is “blind”, and will tell whatever controller is attached to the Request to deal with the command extracted from the Request body, which necessitates the interface polymorphism. UserControllerInterface is also used to track which pages are currently being viewed by users, via maps on Pages.
type UserControllerInterface interface {

	// HandleCommandRegister handles a register command. Only Guest controllers should not automatically produce an error if this is called.
	HandleCommandRegister(*communication.Register, *Server)
	// HandleCommandLogin handles a login command. Only Guest controllers should not respond with an error message if this is called.
	HandleCommandLogin(*communication.Login, *Server)

	// HandleCommandLogout handles a logout command. Guest Controllers should respond with an error message.
	HandleCommandLogout(*communication.Logout, *Server)

	// HandleCommandEmail handles a user's command to change the email associated with their account.
	HandleCommandChangeEmail(*communication.ChangeEmail, *Server)

	// HandleCommandChangeProfileBlurb handles a user's command to change their profile blurb.
	HandleCommandChangeProfileBlurb(*communication.ChangeProfileBlurb, *Server)

	// HandleCommandPasswordResetRequest handles a user's request for a new password by generating a unique code, saving it in the database, and deleting any previous codes for that user.
	HandleCommandPasswordResetRequest(*communication.PasswordReset, *Server)

	// HandleCommandChangePassword handles when a user submits a password reset code. If the code is valid, it's deleted from the database.
	HandleCommandChangePassword(*communication.SetNewPass, *Server)

	// HandleCommandGetComments handles when a user requests the comment data for a particular url.
	HandleCommandGetComments(comm *communication.GetComments, server *Server)

	// HandleCommandCommentReply handles when a user attempts to post a new comment.
	HandleCommandCommentReply(comm *communication.CommentReply, serv *Server)

	// HandleCommandVoteComment handles when a user attempts to vote on a comment.
	HandleCommandVoteComment(comm *communication.CommentVote, serv *Server)

	// HandleCommandViewUsersReport handles the command, only valid for an admin, to view the report on users and usage
	HandleCommandViewUsersReport(comm *communication.ViewUsersReport, serv *Server)

	// HandleCommandViewUsersReport handles the command, only valid for an admin, to view the report on users and usage
	HandleCommandViewFeedback(comm *communication.ViewFeedback, serv *Server)

	// HandleCommandNewFeedback handles a user attempting to post some new feedback for admins to view
	HandleCommandNewFeedback(comm *communication.Feedback, serv *Server)

	// ToggleFeedbackHidden handles the command, only valid for an admin, to toggle whether a given feedback is hidden.
	ToggleFeedbackHidden(comm *communication.ToggleFeedbackHidden, serv *Server)

	// HandleCommandAssignGlobalModerator handles the command, only valid for an admin, to assign or remove a global moderator.
	HandleCommandAssignGlobalModerator(comm *communication.AssignGlobalModerator, serv *Server)

	// HandleCommandAssignAdmin handles the command, only valid for an admin, to grant admin privileges to another user.
	HandleCommandAssignAdmin(comm *communication.AssignAdmin, serv *Server)

	// HandleCommandViewCommentReports handles the command, only valid for a moderator, to view reported comments.
	HandleCommandViewCommentReports(comm *communication.ViewCommentReports, serv *Server)

	// HandleCommandNewReport handles a command, submitted by a user, to submit a comment report.
	HandleCommandNewReport(comm *communication.PostCommentReport, serv *Server)

	// HandleCommandNewReport handles a command, submitted by an admin, to view access logs.
	HandleCommandViewLogs(comm *communication.ViewAccessLogs, serv *Server)

	// Handle a command, valid for a moderator or above, to moderate a comment.
	HandleCommandModerate(comm *communication.Moderate, serv *Server)

	// Handle a command, valid for a moderator or above, to view Moderator Actions
	HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server)

	// GetUser returns the user associated with this controller
	GetUser() *generated.User

	// GetControllerType returns the controller type, eg "GuestController", "AdminController"
	GetControllerType() string

	// GetPage gets the current page associated with this controller.
	GetPage() *Page

	// SetPage sets a controller's page reference to a new page.
	SetPage(page *Page)

	// Respond writes pending server responses to the response writer.
	Respond(w http.ResponseWriter, r *http.Request)
	// SetCookie adds a Token to the pending server responses. (It no longer actually sets a cookie.)
	SetCookie(w http.ResponseWriter, r *http.Request)
	// AddMessage adds a message to the pending server responses.
	AddMessage(success bool, text string)
	// AddMessage adds a ServerResponse of key 'name' containing data given by the data parameter.
	AddWrapped(name string, data interface{})
}

// UserControllerBase provides data members for UserControllers. It does not implement UserControllerInterface fully. Other controllers are defined by extending this Base class and implementing the rest of the interface. Controllers also retain an array of messages that need to be sent to the client, which will be dispatched the next time a request from that user is received
type UserControllerBase struct {
	User             *generated.User
	Page             *Page
	manager          *UserManager
	lastTokenRefresh time.Time
	nextResponse     []interface{}
	// This flag is used when a GuestController logs in or a Member Controller logs out.
	hasloggedin bool
}

// MemberControllerBase provides data members for MemberControllers. It extends UserControllerBase, adding some fields necessary for validation and password reset tracking.
type MemberControllerBase struct {
	UserControllerBase
	canResetPassword bool
}

// GuestController is attached to an HTTP Request Context when a non-logged in user accesses Comment Anywhere.
type GuestController struct {
	UserControllerBase
}

// MemberController is attached to an HTTP Request Context when a regular member accesses Comment Anywhere.
type MemberController struct {
	MemberControllerBase
}

// DomainModeratorController is attached to an HTTP Request Context when a domain moderator accesses Comment Anywhere.
type DomainModeratorController struct {
	MemberControllerBase
	DomainsModerated []string
}

// GlobalModeratorController is attached to an HTTP Request Context when a domain moderator accesses Comment Anywhere.
type GlobalModeratorController struct {
	MemberControllerBase
}

// AdminController is attached to an HTTP Request Context when an administrator accesses Comment Anywhere.
type AdminController struct {
	MemberControllerBase
}
