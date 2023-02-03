package server

import (
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// UserControllerInterface provides method signatures which other UserController types implement. Controller references are attached to HTTP Request Contexts in the first middleware that a Request passes through. Those controller references are subsequently used by API endpoints to execute access-appropriate code associated with a particular user or guest. At the API endpoints, the Server is “blind”, and will tell whatever controller is attached to the Request to deal with the command extracted from the Request body, which necessitates the interface polymorphism. UserControllerInterface is also used to track which pages are currently being viewed by users, via maps on Pages.
type UserControllerInterface interface {

	// HandleCommandRegister handles a register request. Only Guest controllers should not automatically produce an error if this is called.
	HandleCommandRegister(*communication.Register, *Server)

	GetUser() *generated.User
	Respond(w http.ResponseWriter, r *http.Request)
	SetCookie(w http.ResponseWriter, r *http.Request)
	AddMessage(success bool, text string)
}

// UserControllerBase provides data members for UserControllers. It does not implement UserControllerInterface fully. Other controllers are defined by extending this Base class and implementing the rest of the interface. Controllers also retain an array of messages that need to be sent to the client, which will be dispatched the next time a request from that user is received
type UserControllerBase struct {
	User             *generated.User
	manager          *UserManager
	lastTokenRefresh time.Time
	nextResponse     []interface{}
}

// MemberControllerBase provides data members for MemberControllers. It extends UserControllerBase, adding some fields necessary for validation and password reset tracking.
type MemberControllerBase struct {
	UserControllerBase
	canResetPassword bool
}

// GuestController is attached to an HTTP Request Context when a non-logged in user accesses Comment Anywhere.
type GuestController struct {
	UserControllerBase
	hasloggedin bool
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
