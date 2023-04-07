package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// Adds a new Comment Report to the database
func (s *Store) NewCommentReport(comm *communication.PostCommentReport, user int64) (bool, string) {
	cmnt, err := s.Queries.GetCommentByID(context.Background(), comm.CommentID)
	if err != nil {
		return false, "Couldn't find that comment."
	} else {
		if cmnt.Removed.Valid {
			if cmnt.Removed.Bool == true {
				return false, "That comment is already removed."
			}
		}
	}
	params := generated.CreateCommentReportParams{}
	params.Comment = comm.CommentID
	params.ReportingUser = user
	params.Reason.String = comm.Reason
	params.Reason.Valid = true
	err = s.Queries.CreateCommentReport(context.Background(), params)
	if err != nil {
		return false, "Failed to create comment."
	} else {
		return true, "Your report has been submitted."
	}
}

// Gets comment reports for a domain. If domain is the string 'all', it will get all comment reports instead
func (s *Store) GetCommentReportsFor(domain string) (*communication.CommentReports, error) {
	result := &communication.CommentReports{}
	result.Reports = make([]communication.CommentReport, 0, 0)
	if domain == "all" {
		rows, err := s.Queries.GetAllCommentReports(context.Background())
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return result, nil
		} else {
			for _, val := range rows {
				transformed, err := s.transformGeneratedCommentReport1(val)
				if err != nil {
					return nil, err
				}
				result.Reports = append(result.Reports, *transformed)
			}
			return result, nil
		}
	} else {
		d := sql.NullString{
			Valid:  true,
			String: domain,
		}
		rows, err := s.Queries.GetCommentReportsForDomain(context.Background(), d)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return result, nil
		} else {
			for _, val := range rows {
				transformed, err := s.transformGeneratedCommentReport2(val)
				if err != nil {
					return nil, err
				}
				result.Reports = append(result.Reports, *transformed)
			}
			return result, nil
		}
	}
}

// Transforms a generated.CommentReportsRow as returned from 'GetAllCommentReports' into a communication.CommentReport by running the queries needed to populate the data.
func (s *Store) transformGeneratedCommentReport1(item generated.GetAllCommentReportsRow) (*communication.CommentReport, error) {

	ctx := context.Background()
	result := communication.CommentReport{}
	comm, err := s.Queries.GetCommentByID(ctx, item.Comment)
	if err != nil {
		return nil, err
	}
	result.ReportingUserID = item.ReportingUser
	un, err := s.Queries.GetUsername(ctx, item.ReportingUser)
	if err != nil {
		return nil, err
	}
	result.ReportingUsername = un
	err = s.transformGeneratedCommentToCommunicationCommentWithRemoved(&comm, &result.CommentData)
	if err != nil {
		return nil, err
	}
	result.ActionTaken = item.ActionTaken
	result.ReportId = item.ID
	if item.Reason.Valid {
		result.ReasonReported = item.Reason.String
	} else {
		result.ReasonReported = ""
	}
	result.TimeReported = item.TimeCreated.Unix()
	result.Domain = item.Domain.String
	return &result, nil
}

// Transforms a generated.CommentReportsRow as returned from 'GetCommentReportsForDomain' into a communication.CommentReport by running the queries needed to populate the data.
func (s *Store) transformGeneratedCommentReport2(item generated.GetCommentReportsForDomainRow) (*communication.CommentReport, error) {

	ctx := context.Background()
	result := communication.CommentReport{}
	comm, err := s.Queries.GetCommentByID(ctx, item.Comment)
	if err != nil {
		return nil, err
	}
	result.ReportingUserID = item.ReportingUser
	un, err := s.Queries.GetUsername(ctx, item.ReportingUser)
	if err != nil {
		return nil, err
	}
	result.ReportingUsername = un
	err = s.transformGeneratedCommentToCommunicationCommentWithRemoved(&comm, &result.CommentData)
	if err != nil {
		return nil, err
	}
	result.ActionTaken = item.ActionTaken
	result.ReportId = item.ID
	if item.Reason.Valid {
		result.ReasonReported = item.Reason.String
	} else {
		result.ReasonReported = ""
	}
	result.TimeReported = item.TimeCreated.Unix()
	result.Domain = item.Domain.String
	return &result, nil
}
