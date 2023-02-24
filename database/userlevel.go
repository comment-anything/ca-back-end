package database

import (
	"context"
	"database/sql"
	"errors"
)

// IsAdmin returns a bool determining whether a user of given ID is an admin
func (s *Store) IsAdmin(id int64) (bool, error) {
	assnments, err := s.Queries.GetAdminAssignment(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	isAdmin := false
	is_deactivation := false
	for _, a := range assnments {
		is_deactivation = a.IsDeactivation.Valid && a.IsDeactivation.Bool != false
		isAdmin = !is_deactivation
	}
	return isAdmin, nil

}
