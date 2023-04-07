package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/comment-anything/ca-back-end/database/generated"
)

// GetDomainBans gets the list of domains (as strings) that a given user is banned from, if any. If none, it returns an empty string.
func (s *Store) GetDomainBans(id int64) ([]string, error) {
	results, err := s.Queries.GetDomainBans(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s := make([]string, 0, 1)
			return s, nil
		} else {
			return nil, err
		}
	}
	ret := make([]string, len(results))
	for i, v := range results {
		ret[i] = v.BannedFrom
	}
	return ret, nil
}

// AddDomainBan first checks whether target_user_id is an admin or global mod. If they are, nothing in the DB changes and the function returns false with a failure message. Otherwise, the DomainBans table is updated to reflect the new ban status of the user. If the user is already banned, based on the unique PK of user_id+domain in the DomainBans table, it returns a message saying the user is already banned.
func (s *Store) AddDomainBan(target_user_id int64, domain string, banned_by int64, reason *string) (bool, string) {
	is_admin, err := s.IsAdmin(target_user_id)
	if is_admin {
		return false, "You can't ban an admin!"
	}
	is_gmod, err := s.IsGlobalModerator(target_user_id)
	if is_gmod {
		return false, "You can't ban a global moderator!"
	}
	ctx := context.Background()
	p := generated.BanUserFromDomainParams{}
	p.BannedBy = banned_by
	p.BannedFrom = domain
	p.UserID = target_user_id
	err = s.Queries.BanUserFromDomain(ctx, p)
	if err != nil {
		if strings.HasPrefix(err.Error(), "pq: duplicate key value") {
			return false, "That user is already banned."
		}
		return false, err.Error()
	}

	// add the record
	p2 := generated.AddBanRecordParams{}
	p2.Domain.String = domain
	p2.Domain.Valid = true
	p2.Reason.String = *reason
	p2.Reason.Valid = true
	p2.SetBannedTo.Bool = true
	p2.SetBannedTo.Valid = true
	p2.TakenBy = banned_by
	p2.TargetUser = target_user_id
	s.Queries.AddBanRecord(ctx, p2)

	return true, "Banned user from " + domain + "."
}

// RemoveDomainBan deletes a record from the DomainBans table is updated to reflect the new ban status of the user. If the user is already banned, based on the unique PK of user_id+domain in the DomainBans table, it should return some sort of pk already doesnt exist error.
func (s *Store) RemoveDomainBan(target_user_id int64, domain string, unbanned_by int64, reason *string) (bool, string) {
	ctx := context.Background()
	p := generated.UnbanUserFromDomainParams{}
	p.BannedFrom = domain
	p.UserID = target_user_id
	err := s.Queries.UnbanUserFromDomain(ctx, p)
	if err != nil {
		return false, err.Error()
	}
	// add the record
	p2 := generated.AddBanRecordParams{}
	p2.Domain.String = domain
	p2.Domain.Valid = true
	p2.Reason.String = *reason
	p2.Reason.Valid = true
	p2.SetBannedTo.Bool = false
	p2.SetBannedTo.Valid = true
	p2.TakenBy = unbanned_by
	p2.TargetUser = target_user_id
	s.Queries.AddBanRecord(ctx, p2)
	return true, fmt.Sprintf("Unbanned user from %s.", domain)
}

// GlobalBan updates a record in the user's table to indicate that they are now banned. If they were not already banned, it adds a ban record to that table.
func (s *Store) GlobalBan(target_user_id int64, banned_by int64, reason *string) (bool, string) {
	ctx := context.Background()
	status, err := s.Queries.GetUserBanStatus(ctx, target_user_id)
	if status == true {
		return false, "That user is already globally banned."
	}
	err = s.Queries.BanUserGlobally(ctx, target_user_id)
	if err != nil {
		return false, err.Error()
	}
	p2 := generated.AddBanRecordParams{}
	p2.Domain.Valid = false
	p2.Reason.String = *reason
	p2.Reason.Valid = true
	p2.SetBannedTo.Bool = true
	p2.SetBannedTo.Valid = true
	p2.TakenBy = banned_by
	p2.TargetUser = target_user_id
	s.Queries.AddBanRecord(ctx, p2)
	return true, "Globally banned user."
}

// GlobalUnban updates a record in the user's table to indicate that they are now unbanned. If they were not already unbaned, it adds a ban record to that table.
func (s *Store) GlobalUnban(target_user_id int64, banned_by int64, reason *string) (bool, string) {
	ctx := context.Background()
	status, err := s.Queries.GetUserBanStatus(ctx, target_user_id)
	if status == false {
		return false, "That user not currently globally banned."
	}
	err = s.Queries.UnbanUserGlobally(ctx, target_user_id)
	if err != nil {
		return false, err.Error()
	}
	p2 := generated.AddBanRecordParams{}
	p2.Domain.Valid = false
	p2.Reason.String = *reason
	p2.Reason.Valid = true
	p2.SetBannedTo.Bool = false
	p2.SetBannedTo.Valid = true
	p2.TakenBy = banned_by
	p2.TargetUser = target_user_id
	s.Queries.AddBanRecord(ctx, p2)
	return true, "Globally unbanned user."
}
