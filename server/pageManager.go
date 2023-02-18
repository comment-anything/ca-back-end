package server

import (
	"errors"
	"fmt"

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
		user.SetPage(page)
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
		user.Page = newpage
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
			page.LoadComments(comms)
		}

		pm.Pages[pathID] = page
	}
	return &page, nil
}

// Used with the server cli to get some information about the state of PageManager.
func (pm *PageManager) GetPageManagerCountString() string {
	return fmt.Sprintf("Pages Loaded: %d", len(pm.Pages))
}

func (pm *PageManager) UnloadEmptyPage(serv *Server) {}