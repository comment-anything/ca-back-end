package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetProfile will perform the required queries to generate a communication.UserProfile object
func (st *Store) GetCommUser(user *generated.User) (*communication.UserProfile, error) {
	prof := communication.UserProfile{}
	prof.CreatedOn = user.CreatedAt.UnixMilli()

	prof.IsGlobalModerator = false
	if user.ProfileBlurb.Valid {
		prof.ProfileBlurb = user.ProfileBlurb.String
	} else {
		prof.ProfileBlurb = ""
	}
	prof.UserId = user.ID
	prof.Username = user.Username

	isAdmin, _ := st.IsAdmin(user.ID)
	prof.IsAdmin = isAdmin
	return &prof, nil
}

func (st *Store) GetUsername(id int64) string {
	s, _ := st.Queries.GetUsername(context.Background(), id)
	return s
}
