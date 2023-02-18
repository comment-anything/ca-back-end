package database

import (
	"context"
	"fmt"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
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
		err = s.transformGeneratedCommentToCommunicationComment(&val, &ccom)

	}
	return ccomms, nil

}

func (s *Store) transformGeneratedCommentToCommunicationComment(gen *generated.Comment, com *communication.Comment) error {
	com.CommentId = gen.ID
	com.UserId = gen.Author
	com.TimePosted = gen.CreatedAt.Unix()
	if gen.Hidden.Valid {
		com.Hidden = gen.Hidden.Bool
	} else {
		com.Hidden = false
	}
	if gen.Removed.Valid {
		com.Removed = gen.Hidden.Bool
	} else {
		com.Hidden = false
	}
	if gen.Parent.Valid {
		com.Parent = gen.Parent.Int64
	} else {
		com.Parent = 0 // It's a root level comment. Or should this be -1?
	}
	com.Content = gen.Content
	com.Agree.Downs = 0
	com.Agree.Ups = 0
	com.Factual.Downs = 0
	com.Factual.Ups = 0
	com.Funny.Ups = 0
	com.Funny.Downs = 0
	err := s.populateUsername(com)
	if err != nil {
		return err
	}
	err = s.populateVotes(com)
	if err != nil {
		return err
	}
	return nil
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

func (s *Store) NewComment(comm *communication.CommentReply, userId int64, pathId int64) (*communication.Comment, error) {
	result := communication.Comment{}
	ctx := context.Background()
	params := generated.CreateCommentParams{}
	params.Author = userId
	params.Content = comm.Reply
	params.PathID = pathId
	if comm.ReplyingTo != 0 { // 0 is root comment
		fmt.Printf("\nDB.NewComment, we are not replying to root!")
		_, err := s.Queries.GetCommentByID(ctx, comm.ReplyingTo)
		if err != nil {
			fmt.Printf("\nDB.NewComment, parent result: %s", err.Error())
			return nil, err
		}
	}
	params.Parent.Valid = false
	params.Parent.Int64 = comm.ReplyingTo
	gencom, err := s.Queries.CreateComment(ctx, params)
	if err != nil {
		fmt.Printf("\nDB.NewComment, Failure to create comment! %s", err.Error())
		return nil, err
	}
	s.transformGeneratedCommentToCommunicationComment(&gencom, &result)
	return &result, nil

}
