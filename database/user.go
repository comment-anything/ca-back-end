package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	prof.IsAdmin = false
	ad_assnments, err := st.Queries.GetAdminAssignment(context.Background(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			prof.IsAdmin = false
		} else {
			return &prof, errors.New("Error figuring out if admin.")
		}
	}
	for _, a := range ad_assnments {
		fmt.Printf("Checking assignment %v \n\t currently: %v", a, prof.IsAdmin)
		if a.IsDeactivation.Valid && a.IsDeactivation.Bool == prof.IsAdmin {
			prof.IsAdmin = !prof.IsAdmin
		}
		if !a.IsDeactivation.Valid && !prof.IsAdmin {
			prof.IsAdmin = !prof.IsAdmin
		}
	}

	return &prof, nil
}
