package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
)

// GetComments runs queries necessary to get generated.comments for a given path and transform them into an array of communication.Comments
func (s *Store) GetComments(pathID int64) ([]communication.Comment, error) {
	ctx := context.Background()
	rawcomms, err := s.Queries.GetCommentsForPath(ctx, pathID)
	if err != nil {
		return nil, err
	}
	ccomms := make([]communication.Comment, len(rawcomms))
	for i, val := range rawcomms {
		ccom := ccomms[i]
		ccom.CommentId = val.ID
		if val.Hidden.Valid {
			ccom.Hidden = val.Hidden.Bool
		} else {
			ccom.Hidden = false
		}
		if val.Removed.Valid {
			ccom.Removed = val.Hidden.Bool
		} else {
			ccom.Hidden = false
		}
		if val.Parent.Valid {
			ccom.Parent = val.Parent.Int64
		} else {
			ccom.Parent = 0 // It's a root level comment. Or should this be -1?
		}
		ccom.Agree.Downs = 0
		ccom.Agree.Ups = 0
		ccom.Factual.Downs = 0
		ccom.Factual.Ups = 0
		ccom.Funny.Ups = 0
		ccom.Funny.Downs = 0
		err = s.populateUsername(&ccom)
		err = s.populateVotes(&ccom)
		// do we want to do return early if just one comment breaks?

	}
	return ccomms, nil

}

// populateUsername calls database methods needed to populate comm with a user's name
func (s *Store) populateUsername(comm *communication.Comment) error {
	user, err := s.Queries.GetUserByID(context.Background(), comm.UserId)
	if err != nil {
		return err
	}
	comm.Username = user.Username
	return nil
}

// populateVotes calls database methods needed to populate comm with vote data
func (s *Store) populateVotes(comm *communication.Comment) error {
	ctx := context.Background()
	results, err := s.Queries.GetUpVotesForComment(ctx, comm.CommentId)
	if err != nil {
		return err
	}
	for _, val := range results {
		switch val.Category {
		case "funny":
			comm.Funny.Ups += val.Sum
		case "factual":
			comm.Factual.Ups += val.Sum
		case "agree":
			comm.Factual.Ups += val.Sum
		}
	}
	results2, err := s.Queries.GetDownVotesForComment(ctx, comm.CommentId)
	if err != nil {
		return err
	}
	for _, val := range results2 {
		switch val.Category {
		case "funny":
			comm.Funny.Downs += val.Sum
		case "factual":
			comm.Factual.Downs += val.Sum
		case "agree":
			comm.Factual.Downs += val.Sum
		}
	}
	return nil

}
