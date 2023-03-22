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

// CommentVoteRecord contains data for the number of votes on a comment.
type CommentVoteDimension struct {
	AlreadyVoted int8
	Downs        int64
	Ups          int64
}

// Comment provides the data the Front End needs to render a comment.
type Comment struct {
	UserId     int64
	Username   string
	CommentId  int64
	Content    string
	Factual    CommentVoteDimension
	Funny      CommentVoteDimension
	Agree      CommentVoteDimension
	Hidden     bool
	Parent     int64
	Removed    bool
	TimePosted int64
}

// FullPage is returned when a user first requests comments for a new page. It contains an array of all comment data for that page.
type FullPage struct {
	Comments []Comment
	Domain   string
	Path     string
}

// FeedbackRecord contains data the Front End needs to render a FeedbackRecord, which is a record of a user-submitted feedback, viewed by an Admin, such as a feature request, or bug report.
type FeedbackRecord struct {
	Content     string
	Hide        bool
	ID          int64
	SubmittedAt int64
	// May be "bug" | "feature" | "general"
	FeedbackType string
	UserID       int64
	Username     string
}

// FeedbackReport contains an array of feedbackRecords based on the admin's requesting parameters. It is sent to an admin when they request a report on feedbacks
type FeedbackReport struct {
	Records []FeedbackRecord
}

// AdminUsersReport is dispatched when an Admin requests the Users report
type AdminUsersReport struct {
	LoggedInUserCount int64
	NewestUserId      int64
	NewestUsername    string
	UserCount         int64
}

// CommentReport is dispatched when a moderator requests the reported comments
type CommentReport struct {
	ReportId          int64
	ReportingUserID   int64
	ReportingUsername string
	CommentData       Comment
	ReasonReported    string
	ActionTaken       bool
	TimeReported      int64
	Domain            string
}

// CommentReports is dispatched when a moderator requests the reported comments
type CommentReports struct {
	Reports []CommentReport
}

// AdminAccessLogs are dispatched when an admin wants to see what IPs have been accessing the server, which users are associated with them, and what endpoints they are accessing.
type AdminAccessLog struct {
	LogId    int64
	Ip       string
	Url      string
	AtTime   int64
	UserId   int64
	Username string
}

type AdminAccessLogs struct {
	Logs []AdminAccessLog
}
