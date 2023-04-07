package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// ModerateComment will create the relevant moderation record if possible. It does not validate whether the moderator has authority to moderate this comment; such validation is assumed to be already performed in the calling controller.
func (s *Store) ModerateComment(moderator int64, mod *communication.Moderate) (bool, string) {
	ctx := context.Background()
	p := generated.ModerateCommentParams{}
	p.Hidden.Valid = true
	p.Hidden.Bool = mod.SetHiddenTo
	p.Removed.Valid = true
	p.Removed.Bool = mod.SetRemovedTo
	p.ID = mod.CommentID
	err := s.Queries.ModerateComment(ctx, p)
	if err != nil {
		return false, err.Error()
	}
	p2 := generated.CreateModActionRecordParams{}
	if mod.ReportID != nil {
		p2.AssociatedReport.Valid = true
		p2.AssociatedReport.Int64 = *mod.ReportID
		s.Queries.ActionTakenOnReport(ctx, *mod.ReportID)
	}
	p2.CommentID = mod.CommentID
	p2.Reason.Valid = true
	p2.Reason.String = mod.Reason
	p2.SetHiddenTo.Valid = true
	p2.SetHiddenTo.Bool = mod.SetHiddenTo
	p2.SetRemovedTo.Valid = true
	p2.SetRemovedTo.Valid = true
	p2.TakenBy = moderator
	err = s.Queries.CreateModActionRecord(ctx, p2)
	if err != nil {
		return false, err.Error()
	}
	return true, "Comment Moderated"
}
