package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

const (
	CtxLog = string("log")
)

// CreateLoggedRequest should be called at the beginning of an http.Request's lifecycle. It attaches the logID to the request context and returns the request with the new context. It inserts the IP and the API endpoint to the database.
func (s *Store) CreateLoggedRequest(r *http.Request) *http.Request {
	p := generated.CreateLogParams{}
	p.Ip.Valid = true
	p.Ip.String = r.RemoteAddr
	p.Url.Valid = true
	p.Url.String = r.URL.Path

	log, err := s.Queries.CreateLog(context.Background(), p)
	if err == nil {
		newctx := context.WithValue(r.Context(), CtxLog, log.ID)
		return r.WithContext(newctx)
	} else {
		return r
	}
}

// LogUser should be called once a user controller interface is assigned to a request. It updates the database with the user's id.
func (s *Store) LogUser(r *http.Request, id int64) {
	maybe_log := r.Context().Value(CtxLog)
	if maybe_log != nil {
		p := generated.UpdateLogUserParams{}
		p.ID = maybe_log.(int64)
		p.UserID.Valid = true
		p.UserID.Int64 = id
		s.Queries.UpdateLogUser(context.Background(), p)
	}
}

// GetLogs parses a ViewAccessLogs message to determine the filtering, applying the filters to data returned by a database query
func (s *Store) GetLogs(view *communication.ViewAccessLogs) (*communication.AdminAccessLogs, error) {
	c := context.Background()
	p := generated.GetLogsInRangeParams{}
	if view.StartingAt == nil {
		p.AtTime = time.UnixMilli(0)
	} else {
		p.AtTime = time.UnixMilli(*view.StartingAt)
	}
	if view.EndingAt == nil {
		p.AtTime_2 = time.Now()
	} else {
		p.AtTime_2 = time.UnixMilli(*view.EndingAt)
	}
	logs_range, err := s.Queries.GetLogsInRange(c, p)
	if err != nil {
		return nil, err
	}
	if len(view.ForUser) > 0 {
		u_id, err := s.Queries.GetUserByUserName(c, view.ForUser)
		if err != nil {
			return nil, err
		}
		user_filter := make([]generated.Log, 0, len(logs_range))
		for _, val := range logs_range {
			if val.UserID.Valid == true {
				if val.UserID.Int64 == u_id.ID {
					user_filter = append(user_filter, val)
				}
			}
		}
		logs_range = user_filter
	}
	if len(view.ForIp) > 0 {
		ip_filter := make([]generated.Log, 0, len(logs_range))
		for _, val := range logs_range {
			if val.Ip.String == view.ForIp {
				ip_filter = append(ip_filter, val)
			}
		}
		logs_range = ip_filter
	}
	if len(view.ForEndpoint) > 0 {
		endpoint_filter := make([]generated.Log, 0, len(logs_range))
		for _, val := range logs_range {
			fmt.Printf("\n val %s", val.Url.String)
			if val.Url.String == view.ForEndpoint {
				endpoint_filter = append(endpoint_filter, val)
			}
		}
		logs_range = endpoint_filter
	}
	adAcc := communication.AdminAccessLogs{
		Logs: make([]communication.AdminAccessLog, len(logs_range)),
	}
	for i, val := range logs_range {
		adAcc.Logs[i].LogId = val.ID
		adAcc.Logs[i].AtTime = val.AtTime.UnixMilli()
		if val.UserID.Valid {
			adAcc.Logs[i].UserId = val.UserID.Int64
			un, err := s.Queries.GetUserByID(c, val.UserID.Int64)
			if err == nil {
				adAcc.Logs[i].Username = un.Username
			}
		}
		if val.Ip.Valid {
			adAcc.Logs[i].Ip = val.Ip.String
		}
		if val.Url.Valid {
			adAcc.Logs[i].Url = val.Url.String
		}
	}
	return &adAcc, nil

}
