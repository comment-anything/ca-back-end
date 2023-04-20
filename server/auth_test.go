package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/golang-jwt/jwt"
)

func TestGetToken(t *testing.T) {
	os.Clearenv()
	config.Vals.Load("../.env", false)
	defer config.Vals.Reset()
	tstring, err := GetToken(1, true)
	if err != nil {
		t.Errorf("Failed to sign token")
	}

	token, _ := jwt.Parse(tstring, keyfunc)
	claims := token.Claims.(jwt.MapClaims)
	timestring := claims[TokenExpKey].(string)
	timeval, err := time.Parse(time.RFC3339, timestring)
	if err != nil {
		t.Errorf("Time should be able to parse token time: %s", err.Error())
	}
	now := time.Now().Unix()
	if timeval.Unix() <= now {
		t.Errorf("Exp time should be greater than current time, but was %d less", now-timeval.Unix())
	}
	t.Logf("tok expires at %v", timeval)
	id_string := claims[TokenIDKey].(string)
	raw_id_int, err := strconv.Atoi(id_string)
	if err != nil {
		t.Errorf("should be able to gen string from id in token")
	}
	user_id := int64(raw_id_int)
	if user_id != 1 {
		t.Errorf("Token should have same user id after parsing")
	}
	is_guest_string := claims[TokenGuestKey].(string)
	is_guest, err := strconv.ParseBool(is_guest_string)
	if err != nil {
		t.Errorf("Token should parse guest bool")
	}
	if is_guest != true {
		t.Errorf("Token should have correct guest bool val")
	}
}

func TestAuths(t *testing.T) {
	config.Vals.Load("../.env", false)
	defer config.Vals.Reset()
	server, _ := New()
	go server.Start(false)
	time.Sleep(time.Duration(int(time.Second / 10)))
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maybe_controller := r.Context().Value(CtxController)
		if maybe_controller == nil {
			t.Errorf("A guest controller should have been added.")
		}
	})
	rauth := server.EnsureController(nextHandler)
	rauth.ServeHTTP(w, r)

}
