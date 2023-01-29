package communication

// LoginResponse is sent to the client when they successfully log in.
type LoginResponse struct {
	LoggedInAs UserProfile
}

// UserProfile contains data needed by the Front End to display a profile for a user.
type UserProfile struct {
	UserId            int64
	Username          string
	CreatedOn         int64
	DomainsModerating []string
	IsAdmin           bool
	IsDomainModerator bool
	IsGlobalModerator bool
	ProfileBlurb      string
}
