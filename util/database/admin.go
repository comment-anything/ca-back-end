package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/database/generated"
)

func (s *Store) AssignGlobalModerator(to int64, by int64, deactivate bool) error {
	p := generated.CreateGlobalModeratorAssignmentParams{
		AssignedTo: to, AssignedBy: by, IsDeactivation: deactivate,
	}
	err := s.Queries.CreateGlobalModeratorAssignment(context.Background(), p)
	return err
}

func (s *Store) AssignAdmin(to int64, by int64) error {
	p := generated.AssignAdminParams{
		AssignedTo: to, AssignedBy: by,
	}
	err := s.Queries.AssignAdmin(context.Background(), p)
	return err
}
