package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	CtxController = contextKey("controller")
)

const (
	TokenIDKey    = "i"
	TokenExpKey   = "x"
	TokenGuestKey = "g"
)

// keyfunc is used by JWT, it confirms the signing method and returns the secret key for parsing
func keyfunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, errors.New("couldn't parse token: bad signing method of " + token.Method.Alg())
	} // brutal below; have to covert string to byte for signing
	return []byte(config.Vals.Server.JWTKey), nil
}

// GetToken simply returns a JWT token signed with the secret key, with an expiry time of 1 hour and for the userid given as a parameter. It performs no validation.
func GetToken(userid int64, isGuest bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[TokenExpKey] = time.Now().Add(60 * time.Minute).Format(time.RFC3339)
	claims[TokenIDKey] = fmt.Sprint(userid)
	if isGuest {
		claims[TokenGuestKey] = "true"
	} else {
		claims[TokenGuestKey] = "false"
	}
	tstring, err := token.SignedString([]byte(config.Vals.Server.JWTKey))
	if err != nil {
		return "", err
	}
	return tstring, nil
}

// ReadsAuth is a middleware which causes a Controller to be attached to a Request Context if a valid token is present.
func (s *Server) ReadsAuth(handler http.Handler) http.Handler {

	next := handler.ServeHTTP

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cook, err := r.Cookie(config.Vals.Server.JWTCookieName)
		if err != nil {
			next(w, r)
			return
		}
		tstring := cook.Value
		token, err := jwt.Parse(tstring, keyfunc)
		if err != nil {
			next(w, r)
			return
		}
		if !token.Valid {
			next(w, r)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		expires := claims[TokenExpKey].(string)
		t, err := time.Parse(time.RFC3339, expires)
		if err != nil {
			next(w, r)
			return
		}
		expires_at := t.Unix()
		now := time.Now().Unix()
		if now > expires_at {
			next(w, r)
			return
		}
		id_string := claims[TokenIDKey].(string)
		raw_id_int, err := strconv.Atoi(id_string)
		if err != nil {
			next(w, r)
			return
		}
		user_id := int64(raw_id_int)

		is_guest_string := claims[TokenGuestKey].(string)
		is_guest, err := strconv.ParseBool(is_guest_string)
		if err != nil {
			next(w, r)
			return
		}
		controller, err := s.users.GetControllerById(user_id, is_guest)
		if err != nil {
			next(w, r)
			return
		}
		newctx := context.WithValue(r.Context(), CtxController, controller)
		next(w, r.WithContext(newctx))
	})
}

// EnsureController checks if there is a controller attached. If there isn't it attaches a new GuestController.
func (s *Server) EnsureController(handler http.Handler) http.Handler {

	next := handler.ServeHTTP

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maybe_controller := r.Context().Value(CtxController)
		if maybe_controller == nil {
			gc := s.users.CreateGuestController()
			newctx := context.WithValue(r.Context(), CtxController, gc)
			next(w, r.WithContext(newctx))
		} else {
			next(w, r)
		}
	})
}
