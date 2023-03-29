package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandBan tells a guest they do not have permission to perform this action.  */
func (c *GuestController) HandleCommandBan(comm *communication.Ban, serv *Server) {
	c.AddMessage(false, "You do not have permission to ban users.")
}

// AddDomainBan on a guest controller prints an error to the console.
func (c *GuestController) AddDomainBan(s string) {
	fmt.Printf("\nAddDomainBan was called on a guest controller! It shouldn't be ever!")
}

// RemoveDomainBan on a guest controller prints an error to the console.
func (c *GuestController) RemoveDomainBan(s string) {
	fmt.Printf("\nAddDomainBan was called on a guest controller! It shouldn't be, ever!")
}

/** HandleCommandBan tells a member they do not have permission to perform this action.  */
func (c *MemberController) HandleCommandBan(comm *communication.Ban, serv *Server) {
	c.AddMessage(false, "You do not have permission to ban users.")
}

// AddDomainBan on a member controller adds s to the list of domains a member is banned from.
func (c *MemberControllerBase) AddDomainBan(s string) {
	c.BannedFrom = append(c.BannedFrom, s)
}

// RemoveDomainBan on a member controller removes s from the list of domains a member is banned from.
func (c *MemberControllerBase) RemoveDomainBan(s string) {
	narr := make([]string, 0, len(c.BannedFrom))
	for _, v := range c.BannedFrom {
		if v != s {
			narr = append(narr, v)
		}
	}
	c.BannedFrom = narr
}

/** HandleCommandBan for a domain moderator will first check if the domain moderator has authority to ban/unban on this domain. A ban record will be created in the database. If the user is logged in, that domain will be added/removed to/from their list of exclusion domains. */
func (c *DomainModeratorController) HandleCommandBan(comm *communication.Ban, serv *Server) {
	if len(comm.Domain) < 1 {
		c.AddMessage(false, "You do not have permission to ban globally.")
		return
	}
	is_your_domain := false
	for _, v := range c.DomainsModerated {
		if v == comm.Domain {
			is_your_domain = true
			break
		}
	}
	if !is_your_domain {
		c.AddMessage(false, fmt.Sprintf("You do not have permission to moderate %s.", comm.Domain))
		return
	}
	if comm.Ban {
		c.AddMessage(serv.users.DomainBanUser(comm, c.User.ID))
	} else {
		c.AddMessage(serv.users.DomainUnbanUser(comm, c.User.ID))
	}
}

/** HandleCommandBan will attempt to make the necessary calls to ban/unban a user. If globally banned and online, they will be logged out. A ban action record will be created in the database. Domain bans will execute as if they were a domain moderator. GlobalModerators cannot ban other GlobalModerators or admins. */
func (c *GlobalModeratorController) HandleCommandBan(comm *communication.Ban, serv *Server) {
	if len(comm.Domain) > 0 { // Then it's a domain ban.
		if comm.Ban {
			c.AddMessage(serv.users.DomainBanUser(comm, c.User.ID))
		} else {
			c.AddMessage(serv.users.DomainUnbanUser(comm, c.User.ID))
		}
		return
	}
	if comm.Ban == true {
		c.AddMessage(serv.users.GlobalBanUser(comm, c.User.ID))
	} else {
		id, err := serv.DB.GetUserID(comm.Username)
		if err != nil {
			c.AddMessage(false, "Couldn't find that user.")
		} else {
			c.AddMessage(serv.DB.GlobalUnban(id, c.User.ID, &comm.Reason))
		}
	}
}

/** HandleCommandBan for an admin will work as though for a global moderator. To ban a global moderator, an admin must first remove those privileges. */
func (c *AdminController) HandleCommandBan(comm *communication.Ban, serv *Server) {
	if len(comm.Domain) > 0 { // Then it's a domain ban.
		if comm.Ban {
			c.AddMessage(serv.users.DomainBanUser(comm, c.User.ID))
		} else {
			c.AddMessage(serv.users.DomainUnbanUser(comm, c.User.ID))
		}
		return
	}
	if comm.Ban == true {
		c.AddMessage(serv.users.GlobalBanUser(comm, c.User.ID))
	} else {
		id, err := serv.DB.GetUserID(comm.Username)
		if err != nil {
			c.AddMessage(false, "Couldn't find that user.")
		} else {
			c.AddMessage(serv.DB.GlobalUnban(id, c.User.ID, &comm.Reason))
		}
	}
}

// postBan is the API endpoint for when a user attempts to ban or unban a user. It's called when they send a POST request to "/ban". It expects a JSON object of type 'communication.Ban' in the body of the request. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) postBan(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.Ban{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandBan(&comm, s)
		}
	}
}
