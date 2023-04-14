package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

// Gets a random code, saves it in the PW Reset Codes table, then dispatches an email to the user using the util.
func (s *Store) PwResetRequest(comm *communication.PasswordReset) (bool, error) {
	ctx := context.Background()
	user, err := s.Queries.GetUserByEmail(ctx, comm.Email)
	if err != nil {
		// Fake success!
		return true, nil
	}
	err = s.Queries.DeletePreviousPWRestCodesForUser(ctx, user.ID)

	tries := 0
	params := generated.CreatePWResetCodeParams{
		ID:     RandomCode(),
		UserID: user.ID,
	}
	var code generated.PasswordResetCode
	for tries < 10 {
		code, err = s.Queries.CreatePWResetCode(ctx, params)
		if err != nil {
			params.ID = RandomCode()
			tries++
		} else {
			tries = 10
		}
	}
	if err == nil {
		SendPWResetCode(user.Email, code.ID)
	}
	return true, nil

}
