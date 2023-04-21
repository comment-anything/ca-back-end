package database

import (
	"context"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

func (s *Store) VerifyRequest(comm *communication.RequestVerificationCode, email string) (bool, string) {
	ctx := context.Background()
	user, err := s.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "I couldn't find a record of you."
	}
	err = s.Queries.DeletePreviousVerificationCodesForUser(ctx, user.ID)

	tries := 0
	params := generated.CreateVerificationCodeParams{
		ID:     RandomCode(),
		UserID: user.ID,
	}
	var code generated.VerificationCode
	for tries < 10 {
		code, err = s.Queries.CreateVerificationCode(ctx, params)
		if err != nil {
			params.ID = RandomCode()
			tries++
		} else {
			tries = 10
		}
	}
	if err == nil {
		SendVerificationCode(user.Email, user.Username, code.ID)
		return true, "Verification code sent!"
	} else {
		return false, "I couldn't send you a verification code!"
	}
}

func (s *Store) AttemptVerify(comm *communication.InputVerificationCode, userID int64) (bool, string) {
	ctx := context.Background()

	entry, err := s.Queries.GetVerificationCodeEntry(ctx, comm.Code)
	if err != nil {
		return false, "That's not a valid code."
	}
	if entry.UserID != userID {
		return false, "That's not a valid code."
	}
	s.Queries.DeletePreviousVerificationCodesForUser(ctx, userID)
	p := generated.UpdateVerificationParams{
		ID:         userID,
		IsVerified: true,
	}
	s.Queries.UpdateVerification(ctx, p)
	return true, "You have been verified."

}
