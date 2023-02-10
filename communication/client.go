package communication

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

// PasswordResetCode is dispatched by a user when they enter a password reset code. After a user clicks “Forgot My Password”, users may enter the code emailed to them. When they subsequently click the “submit” button, this request is dispatched to the server.
type PasswordResetCode struct {
	Code int32
}

// SetNewPass is dispatched to the Server when the user changes their password. After submitting a valid password reset code, users are prompted to set a new password. When they subsequently click “submit”, this request is dispatched to the server.
type SetNewPass struct {
	Password       string
	RetypePassword string
}
