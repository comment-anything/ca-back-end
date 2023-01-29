package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/config"
)

func TestPostRegister(t *testing.T) {
	config.Vals.Load("../.env")
	defer config.Vals.Reset()
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}
	go server.Start()
	time.Sleep(time.Duration(int(time.Second / 10)))
	req, err := http.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp_recorder := httptest.NewRecorder()
	// call the handler directly
	server.postRegister(resp_recorder, req)
	if status := resp_recorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	/*

		resp.LoggedInAs.Username = "testme"
		resp.LoggedInAs.UserId = 0
		resp.LoggedInAs.IsAdmin = false
		resp.LoggedInAs.CreatedOn = 0
		resp.LoggedInAs.IsDomainModerator = false
		resp.LoggedInAs.IsGlobalModerator = false
		resp.LoggedInAs.ProfileBlurb = "blablablah"
	*/
	actual_response := &communication.LoginResponse{}
	err = json.Unmarshal(resp_recorder.Body.Bytes(), actual_response)
	if err != nil {
		t.Fatalf("The result of the request could not unmarshal: %v", err.Error())
	} else {
		if actual_response.LoggedInAs.Username != "testme" {
			t.Errorf("Incorrect username. Expected %v got %v", "testme", actual_response.LoggedInAs.Username)
		}
	}

	server.Stop()
}
