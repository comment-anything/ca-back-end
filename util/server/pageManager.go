package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/util"
)

// PageManager maintains a map of all instantiated Pages that are currently being viewed by some user or guest. It is responsible for ‘placing’ and ‘removing’ users from pages.
type PageManager struct {
	Pages map[int64]Page
}

// Initializes a new page manager by making the map.
func NewPageManager() PageManager {
	var pages PageManager
	pages.Pages = make(map[int64]Page, 100)
	return pages
}

// MoveMemberToPage is responsible for calling functions on a user's old page to remove the user from that page's map, then calling functions on a new page to add the user to that page's map. In between, LoadPage is called, which may result in a new page instantiation if necessary.
func (pm *PageManager) MoveMemberToPage(user UserControllerInterface, fullpagePath string, serv *Server) {
	page := user.GetPage()
	if page != nil {
		page.RemoveMemberFromPage(user)
	}
	newpage, err := pm.LoadPage(fullpagePath, serv)
	if err != nil {
		fmt.Printf("\n Problem with page load! %s", err.Error())
	} else {
		newpage.AddMemberToPage(user)
	}
}

// MoveGuestToPage is responsible for calling functions on a user's old page to remove the user from that page's map, then calling functions on a new page to add the user to that page's map. In between, LoadPage is called, which may result in a new page instantiation if necessary.
func (pm *PageManager) MoveGuestToPage(user *GuestController, pagePath string, serv *Server) {
	page := user.GetPage()
	if page != nil {
		page.RemoveGuestFromPage(user)
	}
	newpage, err := pm.LoadPage(pagePath, serv)
	if err != nil {
		fmt.Printf("\n Problem with page load! %s", err.Error())
	} else {
		newpage.AddGuestToPage(user)
	}
}

/* LoadPage is called after a GetComments request is received, when a usercontroller moves into a page's map. The domain and subpath are extracted from the raw url provided by the user's communication entity. If a Page does not already exist for the associated url, one will be instantiated, and several other queries will be executed by the store to populate it with comment data. */
func (pm *PageManager) LoadPage(path string, serv *Server) (*Page, error) {
	path_extraction := util.ExtractPathParts(path)
	if path_extraction.Success == false {
		return nil, errors.New(fmt.Sprintf("\nFailed to extract path, got %s %s from %s", path_extraction.Domain, path_extraction.Path, path))
	}

	pathID, err := serv.DB.GetPathResult(path_extraction.Domain, path_extraction.Path)
	if err != nil {
		return nil, err
	}

	page, ok := pm.Pages[pathID]
	if !ok {
		page = NewPage()
		page.pathID = pathID
		page.domain = path_extraction.Domain
		page.path = path_extraction.Path

		comms, err := serv.DB.GetComments(pathID)
		if err != nil {
			fmt.Printf("\nError loading comments for %d: %s", pathID, err.Error())
		}
		page.LoadComments(comms)

		pm.Pages[pathID] = page
	}
	return &page, nil
}

// ModerateComment performs the associated call on the database and returns it to the calling controller. It also calls the relevant page to alter the instanced comment, if such a page exists, so that an instanced comment can be moderated in real-time and the changes can be pushed to any users who may be viewing that page.
func (pm *PageManager) ModerateComment(moderatingUser int64, comm *communication.Moderate, serv *Server) (bool, string) {
	success, msg := serv.DB.ModerateComment(moderatingUser, comm)
	if success {
		commnt, err := serv.DB.Queries.GetCommentByID(context.Background(), comm.CommentID)
		if err == nil {
			page, ok := pm.Pages[commnt.PathID]
			if ok {
				page.ModerateComment(comm, serv)
			}
		}
	}
	return success, msg

}

// Used with the server cli to get some information about the state of PageManager.
func (pm *PageManager) GetPageManagerCountString() string {
	return fmt.Sprintf("Pages Loaded: %d", len(pm.Pages))
}

// gets a list of all the active pages and their ids for the cli.
func (pm *PageManager) GetPagesListString() string {
	var res string
	for key, val := range pm.Pages {
		res += fmt.Sprintf("<%d> : %s | %s\n", key, val.domain, val.path)
	}
	return res
}

// gets info about a page for the cli
func (pm *PageManager) GetPageInfo(s string) string {
	// string like page info <id>
	if len(s) < 10 {
		return "Please supply an id."
	}
	substr := s[10:]
	id, err := strconv.Atoi(substr)
	if err != nil {
		return "Id not understood."
	}
	page, ok := pm.Pages[int64(id)]
	if !ok {
		return "Page not found."
	}
	res := fmt.Sprintf("Page info for %d , domain: %s, path: %s", id, page.domain, page.path)
	res += fmt.Sprintf("\n\t%d guests on page.", len(page.GuestsOnPage))
	res += fmt.Sprintf("\n\t%d comments on page.", len(page.CachedComments))
	for _, mem := range page.MembersOnPage {
		user := mem.GetUser()
		res += fmt.Sprintf("\n\tUser %s on page.", user.Username)
	}
	return res
}

func (pm *PageManager) UnloadEmptyPage(serv *Server) {}
