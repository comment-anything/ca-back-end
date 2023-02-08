package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/comment-anything/ca-back-end/config"
)

type TserverLoggerSettings struct {
	logs *log.Logger
}

// Slogger, or Server Logger, contains the logger object where the server will print when logging is enabled. It may be extended to add additional log options which could be set from the server CLI. If we will implement more options, we should also save them somewhere, or set them in the .env / config object.
var Slogger *TserverLoggerSettings = &TserverLoggerSettings{}

// Returns a middleware function for attaching an ID to an incoming request and logging that request (if logging is enabled). It's called near the beginning of a requests lifecycle.
func LogMiddleware(handler http.Handler) http.Handler {
	Slogger.logs = log.Default()
	next := handler.ServeHTTP

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Vals.Server.DoesLogAll {
			r = logContext(r)
			Slogger.logs.Println(reqString1(r))
		}
		next(w, r)
	})
}

// Returns a middleware function for reading information about a request as it's finishing up and logging it, if logging is enabled. It's called near the end of a request's lifecycle.
func EndLogMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Vals.Server.DoesLogAll {
			Slogger.logs.Println(reqString2(r))
		}
		next(w, r)
	}
}

// The context key for a request ID.
const LogCtxKey string = "RQIdCt"

// The next request ID to use.
var RqId int64 = 1

// Adds a context key with an ID to track the request for logging at the end of the cycle.
func logContext(r *http.Request) *http.Request {
	new_ctx := context.WithValue(r.Context(), LogCtxKey, RqId)
	RqId += 1
	return r.WithContext(new_ctx)
}

// Produces a summary string describing the method and endpoint.
func reqString1(r *http.Request) string {
	id := r.Context().Value(LogCtxKey).(int64)
	return fmt.Sprintf(": [r:%v] %s at %s", id, r.Method, r.URL)
}

// Produces a summary string describing what is being sent back, and to what user.
func reqString2(r *http.Request) string {
	id := r.Context().Value(LogCtxKey).(int64)
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont == nil {
		return fmt.Sprintf(": [r:%v] (done) error: no controller.", id)
	} else {
		user := cont.GetUser()
		return fmt.Sprintf(": [r:%v] (done) responded to %v:%s", id, user.ID, user.Username)
	}
}
