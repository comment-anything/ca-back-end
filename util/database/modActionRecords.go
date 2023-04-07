package database

import (
	"context"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetModRecrods parses a ViewModRecords to determine filtering. It queries the database, applying some filters in the query and some to the returned records, to assembly the final user result.
func (s *Store) GetModRecords(view *communication.ViewModRecords) (*communication.ModRecords, error) {
	c := context.Background()

	// setup query, filter for time or default time range
	p := generated.GetModActionsInRangeParams{}
	if view.From == nil {
		p.TakenOn = time.UnixMilli(0)
	} else {
		p.TakenOn = time.UnixMilli(*view.From)
	}
	if view.To == nil {
		p.TakenOn_2 = time.Now()
	} else {
		p.TakenOn_2 = time.UnixMilli(*view.To)
	}
	actions, err := s.Queries.GetModActionsInRange(c, p)
	if err != nil {
		return nil, err
	}

	// transform generated to communication
	records := communication.ModRecords{}
	records.Records = make([]communication.ModRecord, 0, len(actions))
	for _, v := range actions {
		rec := communication.ModRecord{}
		gencom, err := s.Queries.GetCommentByID(c, v.CommentID)
		if err != nil {
			continue
		}
		err = s.transformGeneratedCommentToCommunicationComment(&gencom, &rec.Comment)
		if err != nil {
			continue
		}
		dom, err := s.Queries.GetCommentDomain(c, gencom.ID)
		if err != nil {
			rec.Domain = dom.String
		}
		rec.ModeratorID = v.TakenBy
		rec.ModeratorUsername = s.GetUsername(v.TakenBy)
		rec.Time = v.TakenOn.Unix()
		if v.AssociatedReport.Valid {
			rec.AssociatedReport = &v.AssociatedReport.Int64
			if v.ReportReason.Valid {
				rec.ReportReason = v.ReportReason.String
			}
			if v.ReportingUser.Valid {
				rec.ReportingUserID = &v.ReportingUser.Int64
				rec.ReportingUsername = s.GetUsername(*rec.ReportingUserID)
			}
			if v.ReportCreated.Valid {
				t := v.ReportCreated.Time.UnixMilli()
				rec.ReportedAt = &t
			}
		}
		if v.ActionReason.Valid {
			rec.Reason = v.ActionReason.String
		}
		if v.SetHiddenTo.Valid {
			rec.SetHiddenTo = &v.SetHiddenTo.Bool
		}
		if v.SetRemovedTo.Valid {
			rec.SetRemovedTo = &v.SetRemovedTo.Bool
		}
		records.Records = append(records.Records, rec)
	}

	// filter by user: will check reporting user, posting user, action-taking user

	if len(view.ByUser) > 0 {
		user_filter := make([]communication.ModRecord, 0, len(records.Records))
		for _, v := range records.Records {
			if v.Comment.Username == view.ByUser || v.ModeratorUsername == view.ByUser || v.ReportingUsername == view.ByUser {
				user_filter = append(user_filter, v)
			}
		}
		records.Records = user_filter
	}

	if len(view.ForDomain) > 0 {
		domain_filter := make([]communication.ModRecord, 0, len(records.Records))
		for _, v := range records.Records {
			d, err := s.Queries.GetCommentDomain(c, v.Comment.CommentId)
			if err != nil && d.String == view.ForDomain {
				domain_filter = append(domain_filter, v)
			}
		}
		records.Records = domain_filter
	}

	// filter by domain

	return &records, nil

}
