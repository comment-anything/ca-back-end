package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

func (s *Store) VoteComment(userID int64, comm *communication.CommentVote) (*communication.Comment, error) {
	ctx := context.Background()
	gencom, err := s.Queries.GetCommentByID(ctx, comm.VotingOn)
	if err != nil {
		return nil, errors.New("Couldn't find the comment you were voting on.")
	}
	if gencom.Removed.Valid {
		if gencom.Removed.Bool == true {
			return nil, errors.New("Can't vote on a removed comment.")
		}
	}
	params := generated.GetVotesForCommentAndCategoryByUserParams{}
	params.CommentID = comm.VotingOn
	params.UserID = userID
	params.Category = comm.VoteType
	vote, err := s.Queries.GetVotesForCommentAndCategoryByUser(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			createParams := generated.CreateCommentVoteParams{}
			createParams.Category = comm.VoteType
			createParams.CommentID = comm.VotingOn
			createParams.UserID = userID
			createParams.Value = comm.Value
			err = s.Queries.CreateCommentVote(ctx, createParams)
			if err != nil {
				return nil, errors.New("Couldn't create vote on that comment.")
			}
		} else {
			return nil, errors.New("Couldn't vote on that comment.")
		}
	}
	newValue := vote.Value + comm.Value
	if newValue == 0 {
		deleteParams := generated.DeleteCommentVoteParams{}
		deleteParams.Category = comm.VoteType
		deleteParams.CommentID = comm.VotingOn
		deleteParams.UserID = userID
		err = s.Queries.DeleteCommentVote(ctx, deleteParams)
		if err != nil {
			return nil, errors.New("Couldn't vote on that comment due to your previous vote.")
		}
	}
	comcom := communication.Comment{}
	s.transformGeneratedCommentToCommunicationCommentWithRemoved(&gencom, &comcom)
	return &comcom, nil

}
