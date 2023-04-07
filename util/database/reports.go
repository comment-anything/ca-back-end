package database

import (
	"context"
	"fmt"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

func (s *Store) GetUserReportDBPartial(c *communication.ViewUsersReport) *communication.AdminUsersReport {
	ctx := context.Background()
	rep := &communication.AdminUsersReport{}
	row, err := s.Queries.GetNewestUser(ctx)
	if err == nil {
		rep.NewestUserId = row.ID
		rep.NewestUsername = row.Username
	}
	count, err := s.Queries.GetUsersCount(ctx)
	if err == nil {
		rep.UserCount = count
	}
	return rep
}

func (s *Store) GetFeedbackReport(c *communication.ViewFeedback) (*communication.FeedbackReport, error) {

	var params generated.GetFeedbacksInRangeOfTypeParams
	params.Hidden = false
	params.SubmittedAt = time.UnixMilli(c.From)
	params.SubmittedAt_2 = time.UnixMilli(c.To)

	var fbacks []generated.Feedback
	var err error

	if c.FeedbackType == "all" {
		var params generated.GetAllFeedbacksInRangeParams
		params.Hidden = false
		params.SubmittedAt = time.UnixMilli(c.From)
		params.SubmittedAt_2 = time.UnixMilli(c.To)
		fbacks, err = s.Queries.GetAllFeedbacksInRange(context.Background(), params)
	} else {
		var params generated.GetFeedbacksInRangeOfTypeParams
		params.SubmittedAt = time.UnixMilli(c.From)
		params.SubmittedAt_2 = time.UnixMilli(c.To)
		params.Type = c.FeedbackType
		fbacks, err = s.Queries.GetFeedbacksInRangeOfType(context.Background(), params)
	}
	if err != nil {
		return nil, err
	}

	rback := &communication.FeedbackReport{}
	rback.Records = make([]communication.FeedbackRecord, 0)
	for _, v := range fbacks {
		rec := &communication.FeedbackRecord{}
		rec.Content = v.Content
		rec.FeedbackType = v.Type
		rec.Hide = v.Hidden
		rec.ID = v.ID
		rec.SubmittedAt = v.SubmittedAt.UnixMilli()
		rec.UserID = v.UserID
		rec.Username = s.GetUsername(v.UserID)
		rback.Records = append(rback.Records, *rec)
	}
	return rback, nil
}

func (s *Store) ToggleFeedback(comm *communication.ToggleFeedbackHidden) (bool, string) {
	var p generated.SetFeedbackHiddenToParams
	p.Hidden = comm.SetHiddenTo
	p.ID = comm.ID
	err := s.Queries.SetFeedbackHiddenTo(context.Background(), p)
	if err != nil {
		return false, err.Error()
	} else {
		var s string
		if p.Hidden {
			s = "hidden"
		} else {
			s = "unhide"
		}
		return true, fmt.Sprintf("Set feedback %d to %s.", p.ID, s)
	}

}
