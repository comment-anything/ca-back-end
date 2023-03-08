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

	if user.ProfileBlurb.Valid {
		prof.ProfileBlurb = user.ProfileBlurb.String
	} else {
		prof.ProfileBlurb = ""
	}
	prof.UserId = user.ID
	prof.Username = user.Username

	isAdmin, _ := st.IsAdmin(user.ID)
	prof.IsAdmin = isAdmin
	if !isAdmin {
		isgmod, err := st.IsGlobalModerator(user.ID)
		if isgmod && err != nil {
			prof.IsGlobalModerator = true
		} else {
			doms_moderated, err := st.GetDomainModeratorAssignments(user.ID)
			if doms_moderated != nil && err != nil {
				prof.IsDomainModerator = true
				prof.DomainsModerating = doms_moderated
			}
		}
	} else {
		prof.IsGlobalModerator = false
		prof.IsDomainModerator = false
	}
	return &prof, nil
}

func (st *Store) GetUsername(id int64) string {
	s, _ := st.Queries.GetUsername(context.Background(), id)
	return s
}
