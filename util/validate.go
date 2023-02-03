package util

import "regexp"

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
	return true, ""
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
