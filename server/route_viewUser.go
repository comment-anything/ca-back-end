package server

import (
	"encoding/json"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

func (u *UserControllerBase) HandleCommandViewUser(comm *communication.ViewUser, serv *Server) {
	puserprof, err := serv.DB.GetPublicUserProfile(comm.Username)
	if err != nil {
		u.AddMessage(false, "Failed to get info for that user.")
		return
	}
	_, logged_in := serv.users.members[puserprof.UserProfile.UserId]
	puserprof.IsLoggedIn = logged_in
	u.AddWrapped("PublicUserProfile", puserprof)
}

func (s *Server) postViewUser(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewUser{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, "I couldn't understand your view user request.")
		} else {
			cont.HandleCommandViewUser(&comm, s)
		}
	}
}
