package util

import "regexp"

// Permits only alphanumerics and '.'
var re_username_chars regexp.Regexp = *regexp.MustCompile(`([a-z]|[A-Z]|[0-9]|\.)`)

// ValidateUsernameCharacters confirms that a username has no prohibited characters that could be used, for example, for sql injection attacks.
func ValidateUsernameCharacters(s string) (bool, string) {
	match := re_username_chars.Match([]byte(s))
	if match {
		return true, ""
	}
	return false, "Usernames may only contain alphanumeric characters."
}

// ValidateUsernameContent confirms a username has no profanity or prohibited characters.
func ValidateUsernameContent(s string) (bool, string) {
	return ValidateUsernameCharacters(s)
}

// ValidatePasswordStrength determines whether a given string is strong enough to be a password.
func ValidatePasswordStrength(s string) (bool, string) {
	if len(s) < 8 {
		return false, "Your password must be at least 8 characters."
	}
	return true, ""
	// insert rainbow table lookup here
}

// ValidateUsername determines whether a given string is acceptable as a username.
func ValidateUsername(s string) (bool, string) {
	if len(s) < 3 {
		return false, "Your username must be at least 3 characters."
	}
	return ValidateUsernameContent(s)
}

// see: https://github.com/google/re2/wiki/Syntax
var re_email regexp.Regexp = *regexp.MustCompile(`[\w\d]+@[\w\d]+\.[\w\d]+`)

// ValidateEmail determines whether a string represents a valid email
func ValidateEmail(s string) (bool, string) {
	if len(s) < 3 {
		return false, "That email isn't long enough to be valid."
	} else {
		match := re_email.Match([]byte(s))
		if !match {
			return false, "That isn't a valid email."
		}
	}
	return true, ""
}

// ValidateProfile runs regexes to determine that a user's profile contains acceptable content.
func ValidateProfile(s string) (bool, string) {
	// TODO: run regexes here
	return true, ""
}
