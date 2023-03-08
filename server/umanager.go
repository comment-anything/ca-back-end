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

// TransferGuest is called when a user logs in. It deletes a previous guest controller from the map, sets the user in the calling guestController to the logged in user, and sets the hasLogged val to true for correct cookie generation.
func (um *UserManager) TransferGuest(oldGuestController *GuestController, user *generated.User) {
	delete(um.guests, oldGuestController.User.ID)

	oldGuestController.User = user
	oldGuestController.hasloggedin = true
}

// TransferMember is called when a user logs out. It deletes a member controller from the map, and sets the calling member controller User to an associated member ID, and sets the hasLoggedIn val to false for correct cookie generation.
func (um *UserManager) TransferMember(oldMemberController *MemberControllerBase) {
	delete(um.members, oldMemberController.User.ID)
	gc := um.CreateGuestController()
	oldMemberController.User.ID = gc.User.ID
	oldMemberController.hasloggedin = false
	delete(um.members, oldMemberController.User.ID)
}

// GetControllerById returns a UserController if it exists in the map and returns an error if it doesnt.
func (um *UserManager) GetControllerById(id int64, isGuest bool) (UserControllerInterface, error) {
	if isGuest {
		cont, ok := um.guests[id]
		if ok {
			return cont, nil
		} else {
			return um.CreateGuestController(), nil
		}
	} else {
		cont, ok := um.members[id]
		if ok {
			fmt.Printf("\nFound controller %d of type %s", id, cont.GetControllerType())
			return cont, nil
		} else {
			newcont, err := um.AttemptCreateMemberController(id)
			fmt.Printf("\nCreated controller %d of type %s", id, newcont.GetControllerType())
			if err != nil {
				return nil, errors.New("Member controller could not be created.")
			}
			return newcont, nil
		}
	}
}

// AttemptCreateMemberController will query the database to see if a member with that id exists. If so, it will add that controller to the map. //TODO: different controllers for various levels of members. Have Admin, need moderators.
func (um *UserManager) AttemptCreateMemberController(id int64) (UserControllerInterface, error) {
	user, err := um.serv.DB.Queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	} else {
		cont, ok := um.members[id]
		if ok {
			return cont, nil // if controller was already made, just return that
		} else {
			// check if they are an admin
			adm, err := um.serv.DB.IsAdmin(id)
			if adm == true && err == nil {
				cont := &AdminController{}
				cont.User = &user
				cont.manager = um
				cont.hasloggedin = true
				um.members[id] = cont
				return cont, nil
			}
			gmod, err := um.serv.DB.IsGlobalModerator(id)
			if gmod == true && err == nil {
				cont := &GlobalModeratorController{}
				cont.User = &user
				cont.manager = um
				cont.hasloggedin = true
				um.members[id] = cont
				return cont, nil
			}
			dmod, err := um.serv.DB.GetDomainModeratorAssignments(id)
			if dmod != nil && err != nil {
				cont := &DomainModeratorController{}
				cont.User = &user
				cont.manager = um
				cont.hasloggedin = true
				cont.DomainsModerated = dmod
				um.members[id] = cont
				return cont, nil
			}
			cont := &MemberController{}
			cont.User = &user
			cont.manager = um
			cont.hasloggedin = true
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
	gc.hasloggedin = false
	um.guests[id] = gc
	return gc
}

func (um *UserManager) GetUserCountString() string {
	return fmt.Sprintf("%v users active; %v members and %v guests.", len(um.members)+len(um.guests), len(um.members), len(um.guests))
}

func (um *UserManager) GetUserListString() string {
	userlist := " Users Online: "
	for _, v := range um.members {
		userlist += v.GetUser().Username + "  "
	}
	return userlist
}
