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
