package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetProfile will perform the required queries to generate a communication.UserProfile object
func (st *Store) GetCommUser(user *generated.User) (*communication.UserProfile, error) {
	prof := &communication.UserProfile{}
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
		if isgmod && err == nil {
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
	prof.IsVerified = user.IsVerified
	bans, err := st.GetDomainBans(user.ID)
	if err == nil {
		prof.DomainsBannedFrom = bans
	}
	return prof, nil
}

// Gets a username given an ID
func (st *Store) GetUsername(id int64) string {
	s, _ := st.Queries.GetUsername(context.Background(), id)
	return s
}

// Gets a user's ID given a user name
func (st *Store) GetUserID(name string) (int64, error) {
	u, err := st.Queries.GetUserByUserName(context.Background(), name)
	if err != nil {
		return 0, err
	} else {
		return u.ID, nil
	}
}

func (st *Store) GetPublicUserProfile(name string) (*communication.PublicUserProfile, error) {
	ctx := context.Background()
	user, err := st.Queries.GetUserByUserName(ctx, name)
	if err != nil {
		return nil, err
	}
	comuser, err := st.GetCommUser(&user)
	if err != nil {
		return nil, err
	}
	retval := &communication.PublicUserProfile{
		UserProfile: *comuser,
		IsLoggedIn:  false,
	}
	return retval, nil
}
