package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
)

/** HandleCommandViewModRecords tells a guest they must be logged in to view mod action records.  */
func (c *GuestController) HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server) {
	c.AddMessage(false, "You do not have permission to view moderator records.")
}

/** HandleCommandViewModRecords tells a guest they must be logged in to view mod action records.  */
func (c *MemberController) HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server) {
	c.AddMessage(false, "You do not have permission to view moderator records.")
}

/** HandleCommandViewModRecords tells a guest they must be logged in to view mod action records.  */
func (c *DomainModeratorController) HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server) {
	good := false
	if len(comm.ForDomain) > 0 {
		for _, v := range c.DomainsModerated {
			if v == comm.ForDomain {
				good = true
			}
		}
	}
	if !good {
		c.AddMessage(false, "You don't have permission to view records for that domain.")
		return
	}
	recs, err := serv.DB.GetModRecords(comm)
	if err != nil {
		c.AddMessage(false, "Failed to get mod records.")
		return
	}
	filt_recs := communication.ModRecords{}
	filt_recs.Records = make([]communication.ModRecord, 0, len(recs.Records))
	for _, v := range recs.Records {
		for _, mydom := range c.DomainsModerated {
			if v.Domain == mydom {
				filt_recs.Records = append(filt_recs.Records, v)
				break
			}
		}
	}
	c.AddWrapped("ModRecords", filt_recs)
}

/** HandleCommandViewModRecords tells a guest they must be logged in to view mod action records.  */
func (c *GlobalModeratorController) HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server) {
	recs, err := serv.DB.GetModRecords(comm)
	if err != nil {
		c.AddMessage(false, err.Error())
	} else {
		c.AddWrapped("ModRecords", recs)
	}
}

/** HandleCommandViewModRecords tells a guest they must be logged in to view mod action records.  */
func (c *AdminController) HandleCommandViewModRecords(comm *communication.ViewModRecords, serv *Server) {
	recs, err := serv.DB.GetModRecords(comm)
	if err != nil {
		c.AddMessage(false, err.Error())
	} else {
		c.AddWrapped("ModRecords", recs)
	}
}

// getModRecords is the API endpoint for when a user attempts to view a feedback report. It's called when they send a POST request to "/ViewModRecords". It expects a JSON object of type 'communication.ViewModRecords'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the response-populating logic.
func (s *Server) getModRecords(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.ViewModRecords{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your request: %s", err.Error()))
		} else {
			cont.HandleCommandViewModRecords(&comm, s)
		}
	}
}
