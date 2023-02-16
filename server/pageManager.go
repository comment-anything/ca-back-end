package server

// PageManager maintains a map of all instantiated Pages that are currently being viewed by some user or guest. It is responsible for ‘placing’ and ‘removing’ users from pages.
type PageManager struct {
	Pages map[string]Page
}

func NewPageManager() PageManager {
	var pages PageManager
	pages.Pages = make(map[string]Page, 100)
	return pages
}

func (pm *PageManager) MoveMemberToPage(user UserControllerInterface, pagePath string, serv *Server) {
	page := user.GetPage()
	if page != nil {
		page.RemoveMemberFromPage(user)
	}
	newpage := pm.LoadPage(pagePath, serv)
	newpage.AddMemberToPage(user)
}

func (pm *PageManager) MoveGuestToPage(user *GuestController, pagePath string, serv *Server) {
	page := user.GetPage()
	if page != nil {
		page.RemoveGuestFromPage(user)
	}
	newpage := pm.LoadPage(pagePath, serv)
	newpage.AddGuestToPage(user)
}

func (pm *PageManager) LoadPage(path string, serv *Server) *Page {
	page, ok := pm.Pages[path]
	if !ok {
		page = NewPage()
		page.fullPath = path
		pm.Pages[path] = page
	}
	return &page
}

func (pm *PageManager) UnloadEmptyPage(serv *Server) {}
