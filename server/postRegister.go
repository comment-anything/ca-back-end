package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// postRegister is the API endpoint for when a user attempts to register a new account. It expects a JSON object of type 'communication.Register'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. PostRegister then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual registration logic.
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
	}
}
