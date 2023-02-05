package server

import (
	"encoding/json"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetUser is just a standard getter for getting the underlying user struct associated with the controller.
func (c *UserControllerBase) GetUser() *generated.User {
	return c.User
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
	}
}

// AddMessage adds an object of type Message wrapped, like every other item in c.nextResponse, in a ServerResponse struct. The success parameter will usually only have the effect of possible changing how a message appears in the front end.
func (c *UserControllerBase) AddMessage(success bool, text string) {
	c.nextResponse = append(c.nextResponse, communication.GetMessage(success, text))
}
