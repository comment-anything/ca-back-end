package server

import (
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/config"
)

// SetCookie writes the Response header "Set-Cookie" with a Name as configured in the .env file and a value of a JWT token encrypted with the key as set in the env file. This Token allows a user to be associated with a particular GuestController, primarily so we can see what resources they are using (e.g., what page they are viewing) and can keep those resources cached for them.
func (c *GuestController) SetCookie(w http.ResponseWriter, r *http.Request) {
	tok, err := GetToken(c.User.ID, !c.hasloggedin)
	if err != nil {
		return
	}
	c.lastTokenRefresh = time.Now()
	auth_cookie := http.Cookie{
		Name:    config.Vals.Server.JWTCookieName,
		Value:   tok,
		MaxAge:  0,
		Path:    "/",
		Expires: c.lastTokenRefresh.Add(60 * time.Minute),
	}
	http.SetCookie(w, &auth_cookie)
}
