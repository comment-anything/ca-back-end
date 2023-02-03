package communication

// ServerResponse wraps server-client communication entities, providing information to the client as to what type of entity data is.
type ServerResponse struct {
	Name string
	Data interface{}
}

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

// Message is a general communication entity used to provide feedback to a client that some action has completed (or not completed) on requests where the client has not asked for any particular data.
type Message struct {
	Success bool
	Text    string
}
