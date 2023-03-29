package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/comment-anything/ca-back-end/communication"
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

// DomainBanUser registers a domain ban for a user. It calls the relevant serv.DB function to add the record in the database. If the user is currently logged in, it adds the domain to their instanced slice of domains they are banned from. If the user is currently logged in, it also pushes a new login response to the user's nextResponse array so that the changes to their profile will be visible to them on the front end.
func (um *UserManager) DomainBanUser(ban *communication.Ban, banned_by int64) (bool, string) {
	user_id, err := um.serv.DB.GetUserID(ban.Username)
	if err != nil {
		return false, fmt.Sprintf("Could not find user %s.", ban.Username)
	}
	cont, ok := um.members[user_id]
	if ok {
		ctype := cont.GetControllerType()
		if ctype == "AdminController" {
			return false, "You cannot ban an admin."
		}
		if ctype == "GlobalModeratorController" {
			return false, "You cannot ban a global moderator."
		}
		ok, why := um.serv.DB.AddDomainBan(user_id, ban.Domain, banned_by, &ban.Reason)
		if !ok {
			return ok, why
		}
		if ctype == "DomainModeratorController" {
			cont.AddDomainBan(ban.Domain)
			cont.AddProfileUpdateResponse()
			return ok, why
		}
		if ctype == "MemberController" {
			cont.AddDomainBan(ban.Domain)
			cont.AddProfileUpdateResponse()
			return ok, why
		}
	}
	return um.serv.DB.AddDomainBan(user_id, ban.Domain, banned_by, &ban.Reason)
}

// DomainBanUser registers a domain unban for a user. It calls the relevant serv.DB function to add the record in the database. If the user is currently logged in, it removes the domain to their instanced slice of domains they are banned from. If the user is currently logged in, it also pushes a new login response to the user's nextResponse array so that the changes to their profile will be visible to them on the front end.
func (um *UserManager) DomainUnbanUser(ban *communication.Ban, unbanned_by int64) (bool, string) {
	user_id, err := um.serv.DB.GetUserID(ban.Username)
	if err != nil {
		return false, fmt.Sprintf("Could not find user %s.", ban.Username)
	}
	cont, ok := um.members[user_id]
	if ok {
		ok, why := um.serv.DB.RemoveDomainBan(user_id, ban.Domain, unbanned_by, &ban.Reason)
		if !ok {
			return ok, why
		}
		cont.RemoveDomainBan(ban.Domain)
		cont.AddProfileUpdateResponse()

	}
	return um.serv.DB.RemoveDomainBan(user_id, ban.Domain, unbanned_by, &ban.Reason)
}

// GlobalBanUsers calls DB functions necessary to realize a user global ban. It also logs the user out, if they are logged in.
func (um *UserManager) GlobalBanUser(ban *communication.Ban, banned_by int64) (bool, string) {
	user_id, err := um.serv.DB.GetUserID(ban.Username)
	if err != nil {
		return false, fmt.Sprintf("Could not find user %s.", ban.Username)
	}
	cont, ok := um.members[user_id]
	if ok {
		ctype := cont.GetControllerType()
		if ctype == "AdminController" {
			return false, "You cannot ban admins."
		}
		if ctype == "GlobalModeratorController" {
			return false, "You cannot ban global moderators."
		}
		user := cont.GetUser()
		user.Banned = true
		return um.serv.DB.GlobalBan(user_id, banned_by, &ban.Reason)
	}
	return um.serv.DB.GlobalBan(user_id, banned_by, &ban.Reason)
}

// TransferMember is called when a user logs out. It deletes a member controller from the map, and sets the calling member controller User to an associated member ID, and sets the hasLoggedIn val to false for correct cookie generation.
func (um *UserManager) TransferMember(oldMemberController *MemberControllerBase) {
	delete(um.members, oldMemberController.User.ID)
	gc := um.CreateGuestController()
	oldMemberController.User.ID = gc.User.ID
	oldMemberController.hasloggedin = false
	delete(um.members, oldMemberController.User.ID)
}

// Swaps an existing member controller to a global mod controller and pushes a login response to their next message to update them on their next action. Used when a logged-in domain moderator or lower is granted global mod privileges.
func (um *UserManager) ChangeMemberControllerToGlobalModController(id int64) (bool, string) {
	mem, ok := um.members[id]
	if !ok {
		return true, "Assignment added; member will have mod privileges on next login."
	}
	curtype := mem.GetControllerType()
	if curtype == "GlobalModeratorController" || curtype == "AdminController" {
		return false, "Member already has privileges."
	}
	curpage := mem.GetPage()
	if curpage != nil {
		curpage.RemoveMemberFromPage(mem)
	}
	delete(um.members, id)
	gmod := &GlobalModeratorController{}
	gmod.User = mem.GetUser()
	gmod.manager = um
	gmod.hasloggedin = true
	gmod.SetPage(curpage)
	um.members[id] = gmod
	gmod.AddProfileUpdateResponse()
	return true, "Assignment added; logged in member given Global Moderator Controller"
}

// Swaps an existing member controller to an admin controller and pushes a login response to their next message to update them on their next action. Used when a logged-in domain moderator or lower is granted admin privileges.
func (um *UserManager) ChangeMemberControllerToAdminController(id int64) (bool, string) {
	mem, ok := um.members[id]
	if !ok {
		return true, "Assignment added; member will have admin privileges on next login."
	}
	curtype := mem.GetControllerType()
	if curtype == "AdminController" {
		return false, "Member already has admin privileges."
	}
	curpage := mem.GetPage()
	if curpage != nil {
		curpage.RemoveMemberFromPage(mem)
	}
	delete(um.members, id)
	cont, err := um.AttemptCreateMemberController(id)
	if err != nil {
		return false, fmt.Sprintf("Failed to create controller: %s", err.Error())
	}

	prof, err := um.serv.DB.GetCommUser(mem.GetUser())
	if err != nil {
		return false, fmt.Sprintf("Failed to get prof: %s", err.Error())
	}
	lr := &communication.LoginResponse{}
	lr.Email = mem.GetUser().Email
	lr.LoggedInAs = *prof
	cont.AddWrapped("LoginResponse", lr)
	return true, "Assignment added; Admin added to user"
}

// Removes GlobalMod privileges from a user. Regenerates a login response so changes are realized on their end.
func (um *UserManager) RemoveGlobalModPrivileges(id int64) (bool, string) {
	mem, ok := um.members[id]
	if !ok {
		return true, "Assignment added; member will not have mod privileges on next login."
	}
	curpage := mem.GetPage()
	if curpage != nil {
		curpage.RemoveMemberFromPage(mem)
	}
	delete(um.members, id)
	cont, err := um.AttemptCreateMemberController(id)
	if err != nil {
		return false, fmt.Sprintf("Failed to create controller: %s", err.Error())
	}
	cont.SetPage(curpage)

	prof, err := um.serv.DB.GetCommUser(mem.GetUser())
	if err != nil {
		return false, fmt.Sprintf("Failed to create controller: %s", err.Error())
	}
	lr := &communication.LoginResponse{}
	lr.Email = mem.GetUser().Email
	lr.LoggedInAs = *prof
	cont.AddWrapped("LoginResponse", lr)
	return true, "Assignment removed; GlobalMod removed from user"
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
			if cont.GetUser().Banned == true {
				delete(um.members, id)
				return nil, errors.New("You have been banned from Comment Anywhere.")
			}
			return cont, nil
		} else {
			newcont, err := um.AttemptCreateMemberController(id)
			fmt.Printf("\nCreated controller %d of type %s", id, newcont.GetControllerType())
			if err != nil {
				return nil, errors.New("Member controller could not be created.")
			}
			if newcont.GetUser().Banned == true {
				delete(um.members, id)
				return nil, errors.New("You have been banned from Comment Anywhere.")
			}
			return newcont, nil
		}
	}
}

/*
* AttemptCreateMemberController will query the database to see if a member with that id exists. If so, it will add that controller to the map. If not the appropriate controller based on the user's privilege levels will be created and added to the map (with the key of the member ID). In either case, the interface will be returned unless there was an error.

If the user is a domain moderator, the list of domains they can moderate is added to their slice.

If the user is below a global moderator and is banned from any domains, the list of domains they are banned from is added to their slice.
*/
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
			bans, err_bans := um.serv.DB.GetDomainBans(id)
			dmod, err := um.serv.DB.GetDomainModeratorAssignments(id)
			if dmod != nil && err == nil {
				cont := &DomainModeratorController{}
				cont.User = &user
				cont.manager = um
				cont.hasloggedin = true
				cont.DomainsModerated = dmod
				um.members[id] = cont
				if err_bans == nil {
					cont.BannedFrom = bans
				}
				return cont, nil
			}
			cont := &MemberController{}
			cont.User = &user
			cont.manager = um
			cont.hasloggedin = true
			if err_bans == nil {
				cont.BannedFrom = bans
			}
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
