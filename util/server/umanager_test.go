package server

import "testing"

func TestUserManager(t *testing.T) {
	um := NewUserManager()
	gc := um.CreateGuestController()
	if gc.User.ID != 1 {
		t.Errorf("A new guest controller in an empty map should have ID 1 but it had: %d", gc.User.ID)
	}
	gc2, err := um.GetControllerById(1, true)
	if err != nil {
		t.Errorf("Should be able to get a guest controller from user manager if it exists but got error: %s", err.Error())
	}
	if gc2.GetUser() != gc.GetUser() {
		t.Error("A retrieved controller should be referencing the same underlying user.")
	}
	_, err = um.GetControllerById(1, false)
	if err == nil {
		t.Errorf("Should produce an error if attempt to get user contrlr that doesnt exist")
	}

	_, err = um.GetControllerById(2, true)
	if err == nil {
		t.Errorf("Should produce an error if guest controller doesnt exist")
	}

}
