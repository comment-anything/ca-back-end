package communication

// ServerResponse wraps server-client communication entities, providing information to the client as to what type of entity data is.
type ServerResponse struct {
	Name string
	Data interface{}
}

// LoginResponse is sent to the client when they successfully log in.
type LoginResponse struct {
	LoggedInAs UserProfile
	Email      string
}

// LogoutResponse is sent to the client when they succesfully log out.
type LogoutResponse struct{}

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

// ProfileUpdateResponse is dispatched to the client when a change to their profile has been realized on the server.
type ProfileUpdateResponse struct {
	LoggedInAs UserProfile
	Email      string
}

// Message is a general communication entity used to provide feedback to a client that some action has completed (or not completed) on requests where the client has not asked for any particular data.
type Message struct {
	Success bool
	Text    string
}

// NewPassResponse is dispatched to the client when they try to use a password reset code to reset their password. It indicates whether the password was reset or not.
type NewPassResponse struct {
	Success bool
	Text    string
}

// Token provides the front end with an authentication key they can use to stay logged in.
type Token struct {
	JWT string
}
