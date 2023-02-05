package server

import (
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/config"
)

// SetCookie writes the Response header "Set-Cookie" with a Name as configured in the .env file and a value of a JWT token encrypted with the key as set in the env file. This Token allows a user to remain logged in, and is refreshed whenever the user interacts with the server. If this token is discarded, the user effectively becomes logged out.
func (c *MemberControllerBase) SetCookie(w http.ResponseWriter, r *http.Request) {
	tok, err := GetToken(c.User.ID, true)
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
