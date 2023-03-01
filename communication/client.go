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
