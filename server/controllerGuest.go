package server

import (
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/config"
)

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
