package database

import (
	"testing"

	"github.com/comment-anything/ca-back-end/config"
)

func TestNew(t *testing.T) {
	config.Vals.Reset()
	nstore, err := New(true)
	if err == nil || nstore != nil {
		t.Errorf("if env vars aren't set, new+connect should produce an error")
	}
	nstore, err = New(false)
	if err != nil || nstore == nil {
		t.Errorf("if new called without connect, should get struct with no errors")
	}
	config.Vals.Load("../.env")
	nstore, err = New(true)
	if err != nil {
		t.Logf("\nThe config string is : %s", config.Vals.DB.ConnectString())
		t.Errorf("\nWith a proper config, store should be able to connect. Ensure the server instance is running and credentials are correct. \nError produced::\t %s", err.Error())
	}
	if nstore != nil {
		nstore.Disconnect()
	}
}

func TestDisconnect(t *testing.T) {
	nstore, err := New(true)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = nstore.Disconnect()
	if err != nil {
		t.Errorf(err.Error())
	}
}
