package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetUser is just a standard getter for getting the underlying user struct associated with the controller.
func (c *UserControllerBase) GetUser() *generated.User {
	return c.User
}

// GetPage gets the page associated with the current user.
func (c *UserControllerBase) GetPage() *Page {
	return c.Page
}

// SetPage sets the user controller's reference to the new page.
func (c *UserControllerBase) SetPage(page *Page) {
	c.Page = page
}

// Respond writes the cached data from c.nextResponse to the responseWriter, and clears it. This should virtually always be the only function in the program which actually writes the body of the http.ResponseWriter.
func (c *UserControllerBase) Respond(w http.ResponseWriter, r *http.Request) {
	outp, err := json.Marshal(c.nextResponse)
	if err != nil {
		w.Write(communication.GetErrMsg(false, "ERROR WRITING RESPONSE"))
	} else if outp == nil {
		w.Write(communication.GetErrMsg(true, ""))
	} else {
		w.Write(outp)
		c.nextResponse = nil
		c.nextResponse = make([]interface{}, 0)

	}
}

// SetCookie writes the Response header "Set-Cookie" with a Name as configured in the .env file and a value of a JWT token encrypted with the key as set in the env file. This Token allows a user to be associated with a particular GuestController, primarily so we can see what resources they are using (e.g., what page they are viewing) and can keep those resources cached for them.
func (c *UserControllerBase) SetCookie(w http.ResponseWriter, r *http.Request) {
	tok, err := GetToken(c.User.ID, !c.hasloggedin)
	if err != nil {
		return
	}
	c.lastTokenRefresh = time.Now()
	var tokenResponse communication.Token
	tokenResponse.JWT = tok
	c.nextResponse = append(c.nextResponse, communication.Wrap("Token", tokenResponse))
}

// AddMessage adds an object of type Message wrapped, like every other item in c.nextResponse, in a ServerResponse struct. The success parameter will usually only have the effect of possible changing how a message appears in the front end.
func (c *UserControllerBase) AddMessage(success bool, text string) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(success, text))
}

// Adds a new message to the next response array. The parameter data should be of a type defined in communication.server so it can be processed correctly on the front end.
func (c *UserControllerBase) AddWrapped(name string, data interface{}) {
	c.nextResponse = append(c.nextResponse, communication.Wrap(name, data))
}

func (c *MemberController) GetControllerType() string {
	return "MemberController"
}

func (c *AdminController) GetControllerType() string {
	return "AdminController"
}

func (c *GuestController) GetControllerType() string {
	return "GuestController"
}

func (c *DomainModeratorController) GetControllerType() string {
	return "DomainModeratorController"
}

func (c *GlobalModeratorController) GetControllerType() string {
	return "GlobalModeratorController"
}

func (c *MemberControllerBase) GetControllerType() string {
	return "MemberControllerBase"
}
