package server

import "testing"

func TestSendMail(t *testing.T) {
	err := SendMail("localhost:25", "comment.anywhere@noreply.go", "Your authentication", "112231", "karlmiller127@gmail.com")
	if err != nil {
		t.Errorf(err.Error())
	}
}
