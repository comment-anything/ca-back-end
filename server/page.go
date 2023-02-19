package server

import (
	"fmt"

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

// NewPage initializes a page object and returns it, creating the maps the page needs.
func NewPage() Page {
	var p Page
	p.CachedComments = make(map[int64]communication.Comment, 50)
	p.GuestsOnPage = make(map[int64]*GuestController, 10)
	p.MembersOnPage = make(map[int64]UserControllerInterface, 10)
	return p
}

// RemoveMemberFromPage removes a member user from this page's map and sets that user's page to nil.
func (p *Page) RemoveMemberFromPage(user UserControllerInterface) {
	user_data := user.GetUser()
	_, ok := p.MembersOnPage[user_data.ID]
	if ok {
		delete(p.MembersOnPage, user_data.ID)
	}
	user.SetPage(nil)

}

// RemoveGuestFromPage removes a guest user from this page's map and sets that user's page to nil.
func (p *Page) RemoveGuestFromPage(user *GuestController) {
	user_data := user.GetUser()
	_, ok := p.GuestsOnPage[user_data.ID]
	if ok {
		delete(p.GuestsOnPage, user_data.ID)
	}
	user.SetPage(nil)
}

// AddMemberToPage adds a member to this page's map and sets that member's page to this page.
func (p *Page) AddMemberToPage(user UserControllerInterface) {
	user_data := user.GetUser()
	p.MembersOnPage[user_data.ID] = user
	user.SetPage(p)
}

// AddGuestToPage adds a guest to this page's map and sets that guest's page to this page.
func (p *Page) AddGuestToPage(user *GuestController) {
	user_data := user.GetUser()
	p.GuestsOnPage[user_data.ID] = user
	user.SetPage(p)
}

// GetComments populates a user's next response with a "FullPage" response consisting of all the comments on this page.
func (p *Page) GetComments(user UserControllerInterface) {

	r := make([]communication.Comment, 0, len(p.CachedComments))

	for _, val := range p.CachedComments {
		r = append(r, val)
	}

	var fp communication.FullPage
	fp.Comments = r
	user.AddWrapped("FullPage", fp)
}

// LoadComments loads the CachedComments map from an array of communication.Comment s.
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
		p.UpdateComment(commResult)
	}
	return true, "Created comment."

}

func (p *Page) UpdateComment(com *communication.Comment) {
	p.CachedComments[com.CommentId] = *com
	for _, gst := range p.GuestsOnPage {
		gst.AddWrapped("Comment", *com)
	}
	for _, mem := range p.MembersOnPage {
		mem.AddWrapped("Comment", *com)
	}
}

// VoteComment effects a user's vote for a comment on this page.
func (p *Page) VoteComment(user UserControllerInterface, comm *communication.CommentVote, serv *Server) (bool, string) {
	if comm.VoteType != "funny" && comm.VoteType != "factual" && comm.VoteType != "agree" {
		return false, fmt.Sprintf("%s is not a valid vote dimension.", comm.VoteType)
	}
	comupdate, err := serv.DB.VoteComment(user.GetUser().ID, comm)
	if err != nil {
		return false, err.Error()
	}
	p.UpdateComment(comupdate)
	return true, "Voted on Comment."
}
