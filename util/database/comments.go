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
		err = s.transformGeneratedCommentToCommunicationCommentWithRemoved(&val, &ccom)
		// WHY DO WE NEED THIS NEXT LINE?! but we do!
		ccomms[i] = ccom

	}
	return ccomms, nil

}

// As the lower transform, but also overwrites the username, ID, and content if the comment has been removed.
func (s *Store) transformGeneratedCommentToCommunicationCommentWithRemoved(gen *generated.Comment, com *communication.Comment) error {
	err := s.transformGeneratedCommentToCommunicationComment(gen, com)
	if err != nil {
		if com.Removed {
			com.Username = ""
			com.Content = "~Removed~"
			com.UserId = 0
		}
		return nil
	}
	return err
}

// transformGeneratedCommentToCommunicationCommentWithRemoved runs queries necessary to transform a raw generated.Comment as received from the database into a communication.Comment used by front ends for comment rendering.
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
		com.Removed = gen.Removed.Bool
	} else {
		com.Removed = false
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
			comm.Agree.Ups += val.Sum
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
			comm.Agree.Downs += val.Sum
		}
	}
	return nil

}

// NewComment inserts a new comment into the database. It returns an error if it couldn't insert the comment. It transforms the comment into a communication.Comment so it can be sent to front ends that need the update.
func (s *Store) NewComment(comm *communication.CommentReply, userId int64, pathId int64) (*communication.Comment, error) {
	result := communication.Comment{}
	ctx := context.Background()
	params := generated.CreateCommentParams{}
	params.Author = userId
	params.Content = comm.Reply
	params.PathID = pathId
	if comm.ReplyingTo != 0 { // 0 is root comment
		fmt.Printf("\nDB.NewComment, we are not replying to root! We are replying to %d", comm.ReplyingTo)
		_, err := s.Queries.GetCommentByID(ctx, comm.ReplyingTo)
		if err != nil {
			fmt.Printf("\nDB.NewComment, parent result: %s", err.Error())
			return nil, err
		} else {
			params.Parent.Valid = true
		}
	}
	params.Parent.Int64 = comm.ReplyingTo
	gencom, err := s.Queries.CreateComment(ctx, params)
	if err != nil {
		fmt.Printf("\nDB.NewComment, Failure to create comment! %s Had params: %v", err.Error(), params)
		return nil, err
	}
	s.transformGeneratedCommentToCommunicationCommentWithRemoved(&gencom, &result)
	return &result, nil

}

// Gets the Domain associated with this comment.
func (s *Store) GetCommentDomain(commentID int64) {

}
