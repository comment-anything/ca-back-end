package server

import (
	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetProfile will perform the required queries to generate a communication.UserProfile object
func (s *Server) GetProfile(user *generated.User) communication.UserProfile {
	prof := communication.UserProfile{}
	prof.CreatedOn = user.CreatedAt.UnixMilli()
	prof.IsAdmin = false
	prof.IsGlobalModerator = false
	if user.ProfileBlurb.Valid {
		prof.ProfileBlurb = user.ProfileBlurb.String
	} else {
		prof.ProfileBlurb = ""
	}
	prof.UserId = user.ID
	prof.Username = user.Username
	return prof
}
