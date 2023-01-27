/*
NOTES:

[This stack overflow answer](https://stackoverflow.com/questions/42474259/golang-how-to-live-test-an-http-server) may be useful when thinking about how to test API end points.
*/
package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, error := New()
	if error != nil {
		t.Errorf("Error with the server: %s", error)
	}
}

func BenchMarkNew(b *testing.B) {
	server, error := New()
	_ = server
	_ = error
}
