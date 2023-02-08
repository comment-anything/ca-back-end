package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/comment-anything/ca-back-end/database/generated"
)

type UserManager struct {
	serv    *Server
	members map[int64]UserControllerInterface
	guests  map[int64]UserControllerInterface
}

// NewUserManager returns a new usermanager. Member 'serv' should be manually assigned to the server after retrieving the manager.
func NewUserManager() UserManager {
	var um UserManager
	um.members = make(map[int64]UserControllerInterface, 100)
	um.guests = make(map[int64]UserControllerInterface, 100)
	return um
}

// TransferGuest deletes a previous guest controller from the map, sets the user in the calling guestController to the logged in user, and sets the hasLogged val to true for correct cookie generation.
func (um *UserManager) TransferGuest(oldGuestController *GuestController, user *generated.User) {
	delete(um.guests, oldGuestController.User.ID)
	oldGuestController.User = user
	oldGuestController.hasloggedin = true
}

// GetControllerById returns a UserController if it exists in the map and returns an error if it doesnt.
func (um *UserManager) GetControllerById(id int64, isGuest bool) (UserControllerInterface, error) {
	if isGuest {
		cont, ok := um.guests[id]
		if ok {
			return cont, nil
		} else {
			return nil, errors.New("No guest controller with that id is present.")
		}
	} else {
		cont, ok := um.members[id]
		if ok {
			return cont, nil
		} else {

			return nil, errors.New("No member controller with that id is present.")
		}
	}
}

// AttemptCreateMemberController will query the database to see if a member with that id exists. If so, it will add that controller to the map. //TODO: different controllers for various levels of members
func (um *UserManager) AttemptCreateMemberController(id int64) (UserControllerInterface, error) {
	user, err := um.serv.DB.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	} else {
		cont, ok := um.members[id]
		if ok {
			return cont, nil // if controller was already made, just return that
		} else {
			cont := &MemberController{}
			cont.User = &user
			um.members[id] = cont
			return cont, nil
		}

	}
}

// CreateGuestController initializes a new guest controller with a unique ID and stores it in the map, and returns it.
func (um *UserManager) CreateGuestController() *GuestController {
	gc := &GuestController{}
	gc.User = &generated.User{}
	gc.manager = um
	exists := true
	var id int64 = 0
	for exists {
		_, exists = um.guests[id]
		id++
	}
	gc.User.ID = id
	gc.User.Username = "--Guest--"
	um.guests[id] = gc
	return gc
}

func (um *UserManager) GetUserCountString() string {
	return fmt.Sprintf("%v users active; %v members and %v guests.", len(um.members)+len(um.guests), len(um.members), len(um.guests))
}
