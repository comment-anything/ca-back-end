package server

import "github.com/comment-anything/ca-back-end/communication"

// Page contains cached data for a page, which is a discrete set of comments associated with a particular URL. It also contains a map of all users and guests on the current page.
type Page struct {
	fullPath       string
	CachedComments map[int64]communication.Comment
	MembersOnPage  map[int64]UserControllerInterface
	GuestsOnPage   map[int64]*GuestController
}

func NewPage() Page {
	var p Page
	p.CachedComments = make(map[int64]communication.Comment, 50)
	p.GuestsOnPage = make(map[int64]*GuestController, 10)
	p.MembersOnPage = make(map[int64]UserControllerInterface, 10)
	return p
}

func (p *Page) RemoveMemberFromPage(user UserControllerInterface) {
	user_data := user.GetUser()
	_, ok := p.MembersOnPage[user_data.ID]
	if ok {
		delete(p.MembersOnPage, user_data.ID)
	}

}
func (p *Page) RemoveGuestFromPage(user *GuestController) {
	user_data := user.GetUser()
	_, ok := p.GuestsOnPage[user_data.ID]
	if ok {
		delete(p.GuestsOnPage, user_data.ID)
	}
}
func (p *Page) AddMemberToPage(user UserControllerInterface) {
	user_data := user.GetUser()
	p.MembersOnPage[user_data.ID] = user
}
func (p *Page) AddGuestToPage(user *GuestController) {
	user_data := user.GetUser()
	p.GuestsOnPage[user_data.ID] = user
}

func (p *Page) GetComments(user UserControllerInterface) {

	r := make([]communication.Comment, len(p.CachedComments))

	for _, val := range p.CachedComments {
		r = append(r, val)
	}

	var fp communication.FullPage
	fp.Comments = r
	user.AddWrapped("FullPage", fp)
}
