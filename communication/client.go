package communication

// ChangeProfileBlurb is dispatched to the server when a client updates their profile blurb.
type ChangeProfileBlurb struct {
	NewBlurb string
}

// ChangeEmail is dispatched to the server when a client wants to change their email. They must supply the correct password as well.
type ChangeEmail struct {
	NewEmail string
	Password string
}

// Register is dispatched to the server when the client clicks “Submit” on the register form.
type Register struct {
	Username       string
	Password       string
	RetypePassword string
	Email          string
	AgreedToTerms  bool
}

// Login is dispatched to the server when the client clicks “Submit” on the login form.
type Login struct {
	Username string
	Password string
}

// Logout is dispatched to the server when the client clicks “Logout”. It does not carry any additional data.
type Logout struct{}

// PasswordReset is dispatched to the server when a password reset is requested. The client supplies the email associated with their account.
type PasswordReset struct {
	Email string
}

// SetNewPass is dispatched to the Server when the user changes their password. After submitting a valid password reset code, users are prompted to set a new password. When they subsequently click “submit”, this request is dispatched to the server.
type SetNewPass struct {
	Email          string
	Code           int64
	Password       string
	RetypePassword string
}

// CommentReply is dispatched to the server when a logged-in user submits a reply to an existing comment or posts a new root-level comment on a page.
type CommentReply struct {
	ReplyingTo int64
	Reply      string
}

// CommentVote is dispatched to the server when a logged-in user votes on a comment.
type CommentVote struct {
	VotingOn int64
	VoteType string
	Value    int16
}

// GetComments is dispatched to the server when a user opens the Browser Extension or when they navigate to a new page with the browser extension over. It is a request for all comments associated with the given url.
type GetComments struct {
	Url           string
	SortedBy      string
	SortAscending bool
}

// ViewUsersReport is dispatched to the server when an admin requests a report on the overall users of comment anywhere.
type ViewUsersReport struct {
}

// ViewFeedback is dispatched to the Server when an admin wishes to view feedback submitted by users of Comment Anywhere
type ViewFeedback struct {
	From int64
	To   int64
	// May be "bug" | "feature" | "general" | "all"
	FeedbackType string
}

// ToggleFeedbackHidden is dispatched to the server when an admin wishes to toggle whether a particular feedback entry is hidden and should be shown on future feedback reports. */
type ToggleFeedbackHidden struct {
	ID          int64
	SetHiddenTo bool
}

// Feedback is dispatched to the Server when a user submits feedback, such as a feature request or bug report, on Comment Anywhere */
type Feedback struct {
	FeedbackType string
	Content      string
}

// AssignGlobalModerator is dispatched to the server when an admin grants or removes global moderator privileges to another user.
type AssignGlobalModerator struct {
	User           string
	IsDeactivation bool
}

// AssignAdmin is dispatched to the server when an admin grants another user admin privileges.
type AssignAdmin struct {
	User string
}

// ViewCommentReports is dispatched to the server when a user requests a list of reported comments
type ViewCommentReports struct {
	Domain string
}

// PostCommentReport is dispatched to the server when a user reports a comment.
type PostCommentReport struct {
	CommentID int64
	Reason    string
}

// ViewAccessLogs is dispatched to the server when an admin requests access logs.
type ViewAccessLogs struct {
	ForUser     string
	ForIp       string
	ForEndpoint string
	// This is only a pointer so "nil" can be read when unmarshaling
	StartingAt *int64
	// This is only a pointer so "nil" can be read when unmarshaling
	EndingAt *int64
}

// Moderate is dispatched to the server when a moderator or admin takes a moderation action on a comment.
type Moderate struct {
	ReportID     *int64
	CommentID    int64
	SetHiddenTo  bool
	SetRemovedTo bool
	Reason       string
}

// ViewModRecords is dispatched to the served when an admin requests moderation records.
type ViewModRecords struct {
	From      int64
	To        int64
	ByUser    string
	ForDomain string
}
