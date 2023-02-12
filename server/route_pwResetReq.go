package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/comment-anything/ca-back-end/communication"
	"github.com/comment-anything/ca-back-end/database/generated"
)

func randomCode() int64 {
	return rand.Int63n(4294967295)
}

// HandleCommandPasswordResetRequest handles a user's request for a new password by generating a unique code, saving it in the database, and deleting any previous codes for that user.
func (c *UserControllerBase) HandleCommandPasswordResetRequest(comm *communication.PasswordReset, serv *Server) {
	user, err := serv.DB.Queries.GetUserByEmail(context.Background(), comm.Email)
	if err == nil {
		// delete all other code entries associated with this user; only the most recent will be valid.
		err := serv.DB.Queries.DeletePreviousPWRestCodesForUser(context.Background(), user.ID)

		tries := 0
		params := generated.CreatePWResetCodeParams{}
		params.ID = randomCode()
		params.UserID = user.ID

		var code generated.PasswordResetCode
		for tries < 10 { // very unlikely the code will be already used, but try a few times just in case
			code, err = serv.DB.Queries.CreatePWResetCode(context.Background(), params)
			if err != nil {
				params.ID = randomCode()
				tries++
			} else {
				tries = 10
			}
		}
		// if there is still an error, 10 tries didn't work and its a database problem not a code collision problem.
		if err != nil {
			_ = code // TODO: actually send the email
		}
	}
	// Even if we don't send an email, we pretend we did for information hiding.
	c.AddMessage(true, "A password reset code has been email to you.")
}

// postPasswordResetRequest is the API endpoint for when a user submits a password reset request. It expects a JSON object of type 'communication.PasswordReset'. As with all endpoints, it first extracts the controller that was attached to the request by earlier middleware. It then then decodes the body of the HTTP Request into an expected communnication entity. It passes that entity to the Controller to perform the actual logic.
func (s *Server) postPasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	cont := r.Context().Value(CtxController).(UserControllerInterface)
	if cont != nil {
		comm := communication.PasswordReset{}
		err := json.NewDecoder(r.Body).Decode(&comm)
		if err != nil {
			cont.AddMessage(false, fmt.Sprintf("I couldn't understand your password reset request: %s", err.Error()))
		} else {
			cont.HandleCommandPasswordResetRequest(&comm, s)
		}
	}
}
