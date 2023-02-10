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
