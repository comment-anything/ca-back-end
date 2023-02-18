package server

import (
	"github.com/comment-anything/ca-back-end/communication"
)

// Page contains cached data for a page, which is a discrete set of comments associated with a particular URL. It also contains a map of all users and guests on the current page.
type Page struct {
	domain         string
	path           string
	pathID         int64
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
	user.SetPage(nil)

}
func (p *Page) RemoveGuestFromPage(user *GuestController) {
	user_data := user.GetUser()
	_, ok := p.GuestsOnPage[user_data.ID]
	if ok {
		delete(p.GuestsOnPage, user_data.ID)
	}
	user.SetPage(nil)
}
func (p *Page) AddMemberToPage(user UserControllerInterface) {
	user_data := user.GetUser()
	p.MembersOnPage[user_data.ID] = user
	user.SetPage(p)
}
func (p *Page) AddGuestToPage(user *GuestController) {
	user_data := user.GetUser()
	p.GuestsOnPage[user_data.ID] = user
	user.SetPage(p)
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

func (p *Page) LoadComments(comments []communication.Comment) {
	for _, val := range comments {
		p.CachedComments[val.CommentId] = val
	}
}

/*
* NewComment posts a new user's comment into the database. It returns a bool representing whether adding the comment was succesful and a string for an error message that will be shown to the user if it was not.

MAYBE: It also notifies all users on the page that there is a new comment by adding it to their pending messages

	(could be prone to errors... what if another user's next request is for getting comments on some other page? )
	- hold off on this implementation for now; we will just use GetComments again on the user that posted a comment so at least theirs is updated...
*/
func (p *Page) NewComment(user UserControllerInterface, comm *communication.CommentReply, serv *Server) (bool, string) {
	commResult, err := serv.DB.NewComment(comm, user.GetUser().ID, p.pathID)
	if err != nil {
		return false, "Couldn't create the comment."
	} else {
		p.CachedComments[commResult.CommentId] = *commResult
		for _, gst := range p.GuestsOnPage {
			gst.AddWrapped("Comment", *commResult)
		}
		for _, mem := range p.MembersOnPage {
			mem.AddWrapped("Comment", *commResult)
		}
	}
	return true, "Created comment."

}
