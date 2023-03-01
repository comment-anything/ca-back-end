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
	if len(assnments) > 0 {
		last_row := assnments[len(assnments)-1]
		last_row_indicates_admin := !last_row.IsDeactivation.Valid || last_row.IsDeactivation.Bool == false
		return last_row_indicates_admin, nil
	} else {
		return false, nil
	}

}
