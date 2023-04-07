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

// GetDomainModeratorAssignments determines which domains a user is permitted to moderate. It returns that list as a slice. If the user is not assigned any domains, nil is returned. The error indicates if there was a problem connecting the database.
func (s *Store) GetDomainModeratorAssignments(userID int64) ([]string, error) {
	results, err := s.Queries.GetDomainModeratorAssignments(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	retval := make([]string, 0, 5)
	for _, v := range results {
		if v.IsDeactivation != true {
			retval = append(retval, v.Domain)
		}
	}
	if len(retval) > 0 {
		return retval, nil
	} else {
		return nil, nil
	}
}

// IsGlobalModerator determines whether a user has been assigned a global moderator role. It returns an error if there was an issue connecting to the database.
func (s *Store) IsGlobalModerator(userID int64) (bool, error) {
	assnments, err := s.Queries.GetGlobalModeratorAssignments(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	if len(assnments) > 0 {
		first_row := assnments[0]
		return first_row.IsDeactivation == false, nil // return true if they weren't most recently deactivated
	} else {
		return false, nil
	}

}
