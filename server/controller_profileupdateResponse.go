package server

import (
	"fmt"

	"github.com/comment-anything/ca-back-end/communication"
)

// AddLoginResponse builds a communication.LoginResponse and fills it with what data is available to a MemberController. This is used for pushing user profile update responses.
func (uc *MemberController) AddProfileUpdateResponse() {
	lr := communication.LoginResponse{}
	lr.Email = uc.User.Email
	lr.LoggedInAs.CreatedOn = uc.User.CreatedAt.Unix()
	lr.LoggedInAs.DomainsBannedFrom = uc.BannedFrom
	lr.LoggedInAs.ProfileBlurb = uc.User.ProfileBlurb.String
	lr.LoggedInAs.Username = uc.User.Username
	lr.LoggedInAs.IsVerified = uc.User.IsVerified
	lr.LoggedInAs.IsAdmin = false
	lr.LoggedInAs.IsDomainModerator = false
	lr.LoggedInAs.IsGlobalModerator = false
	uc.AddWrapped("ProfileUpdateResponse", lr)
}

// AddLoginResponse builds a communication.LoginResponse and fills it with what data is available to a DomainModeratorController. This is used for pushing user profile update responses.
func (uc *DomainModeratorController) AddProfileUpdateResponse() {
	lr := communication.LoginResponse{}
	lr.Email = uc.User.Email
	lr.LoggedInAs.CreatedOn = uc.User.CreatedAt.Unix()
	lr.LoggedInAs.DomainsBannedFrom = uc.BannedFrom
	lr.LoggedInAs.ProfileBlurb = uc.User.ProfileBlurb.String
	lr.LoggedInAs.Username = uc.User.Username
	lr.LoggedInAs.DomainsModerating = uc.DomainsModerated
	lr.LoggedInAs.IsVerified = uc.User.IsVerified
	lr.LoggedInAs.IsAdmin = false
	lr.LoggedInAs.IsDomainModerator = true
	lr.LoggedInAs.IsGlobalModerator = false
	uc.AddWrapped("ProfileUpdateResponse", lr)
}

// AddLoginResponse builds a communication.LoginResponse and fills it with what data is available to a GlobalModeratorController. This is used for pushing user profile update responses.
func (uc *GlobalModeratorController) AddProfileUpdateResponse() {
	lr := communication.LoginResponse{}
	lr.Email = uc.User.Email
	lr.LoggedInAs.CreatedOn = uc.User.CreatedAt.Unix()
	lr.LoggedInAs.DomainsBannedFrom = uc.BannedFrom
	lr.LoggedInAs.ProfileBlurb = uc.User.ProfileBlurb.String
	lr.LoggedInAs.IsVerified = uc.User.IsVerified
	lr.LoggedInAs.Username = uc.User.Username
	lr.LoggedInAs.DomainsModerating = nil
	lr.LoggedInAs.IsAdmin = false
	lr.LoggedInAs.IsDomainModerator = false
	lr.LoggedInAs.IsGlobalModerator = true
	uc.AddWrapped("ProfileUpdateResponse", lr)
}

// AddLoginResponse builds a communication.LoginResponse and fills it with what data is available to an AdminController. This is used for pushing user profile update responses.
func (uc *AdminController) AddProfileUpdateResponse() {
	lr := communication.LoginResponse{}
	lr.Email = uc.User.Email
	lr.LoggedInAs.CreatedOn = uc.User.CreatedAt.Unix()
	lr.LoggedInAs.DomainsBannedFrom = uc.BannedFrom
	lr.LoggedInAs.ProfileBlurb = uc.User.ProfileBlurb.String
	lr.LoggedInAs.IsVerified = uc.User.IsVerified
	lr.LoggedInAs.Username = uc.User.Username
	lr.LoggedInAs.DomainsModerating = nil
	lr.LoggedInAs.IsAdmin = true
	lr.LoggedInAs.IsDomainModerator = false
	lr.LoggedInAs.IsGlobalModerator = false
	uc.AddWrapped("ProfileUpdateResponse", lr)
}

// AddLoginResponse on a guest controller prints an error message. There is no circumstance this should be called, but GuestController needs to implement to be polymorphic with the other controllers.
func (uc *GuestController) AddProfileUpdateResponse() {
	fmt.Printf("\nError: GuestController tried to push a profile update response!")
}
