package server

import (
	"encoding/json"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

// postRegister is the endpoint for when a user attempts to register a new account. It expects a JSON object of type 'Register' and will ultimately return an object of type 'LoginResponse'.
func (s *Server) postRegister(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		cont.AddMessage(true, "TEST MESG")
		//cont.HandleCommandRegister()
	} else {
		resp := &communication.LoginResponse{}
		resp.LoggedInAs.Username = "testme"
		resp.LoggedInAs.UserId = 0
		resp.LoggedInAs.IsAdmin = false
		resp.LoggedInAs.CreatedOn = 0
		resp.LoggedInAs.IsDomainModerator = false
		resp.LoggedInAs.IsGlobalModerator = false
		resp.LoggedInAs.ProfileBlurb = "blablablah"

		respString, _ := json.Marshal(resp)
		w.Write(respString)
	}

}
