package server

import (
	"encoding/json"
	"net/http"

	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetUser gets the underlying user struct associated with the controller.
func (c *UserControllerBase) GetUser() *generated.User {
	return c.User
}

// Respond writes the cached data from c.nextResponse to the responseWriter, and clears it.
func (c *UserControllerBase) Respond(w http.ResponseWriter, r *http.Request) {
	outp, err := json.Marshal(c.nextResponse)
	if err != nil {
		w.Write([]byte("ERROR WRITING RESPONSE"))
	} else {
		w.Write(outp)
	}
}
