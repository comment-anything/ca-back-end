package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// postRegister is the endpoint for when a user attempts to register a new account. It expects a JSON object of type 'Register' and will ultimately return an object of type 'LoginResponse'.
func (s *Server) postRegister(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Register{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your registration request: %s", err.Error()))
		} else {
			cont.HandleCommandRegister(&comm, s)
		}

		//cont.HandleCommandRegister()
	}

}
