package database

import (
	"context"
	"database/sql"

	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetPathResult attempts to get a path with the specified domain path. If no such path is found in the database, it creates a new one. It returns the ID of the path.
func (s *Store) GetPathResult(domain string, path string) (int64, error) {
	ctx := context.Background()
	var params generated.GetPathParams
	params.Domain.Valid = true
	params.Domain.String = domain
	params.Path = sql.NullString{Valid: true, String: path}
	pathOb, err := s.Queries.GetPath(ctx, params)
	if err != nil {
		var cprams generated.CreatePathParams
		s.Queries.EnsureDomainRecordExits(ctx, domain)
		cprams.Domain.Valid = true
		cprams.Domain.String = domain
		cprams.Path = sql.NullString{Valid: true, String: path}
		pathOb, err = s.Queries.CreatePath(ctx, cprams)
	}
	return pathOb.ID, err
}
