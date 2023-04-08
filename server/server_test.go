/*
NOTES:

[This stack overflow answer](https://stackoverflow.com/questions/42474259/golang-how-to-live-test-an-http-server) may be useful when thinking about how to test API end points.
*/
package server

import (
	"testing"
	"time"

	"github.com/comment-anything/ca-back-end/config"
)

func TestNew(t *testing.T) {
	_, error := New()
	if error == nil {
		t.Errorf("The server should return an error if the config file isn't loaded.")
	}
	config.Vals.Load("../.env", false)
	_, error = New()
	if error != nil {
		t.Errorf("Error with the server: %s", error)
	}
	config.Vals.Reset()
}

func TestSetupRouter(t *testing.T) {
	s := &Server{}
	s.setupRouter()
	if s.httpServer.Handler == nil {
		t.Errorf("Server did not set up router.")
	}
}

func shutdown(s *Server, secs float32) {
	time.Sleep(time.Duration(secs * float32(time.Second)))
	s.Stop()
}

func TestStart(t *testing.T) {
	config.Vals.Load("../.env", false)
	defer config.Vals.Reset()
	s, _ := New()
	go shutdown(s, 0.25)
	s.Start(false)
}

func TestStop(t *testing.T) {
	config.Vals.Load("../.env", false)
	defer config.Vals.Reset()
	s, _ := New()
	go s.Start(false)
	time.Sleep(2)
	s.Stop()
}
