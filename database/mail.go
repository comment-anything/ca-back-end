package database

import (
	"fmt"
	"math/rand"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/go-mail/mail"
)

func RandomCode() int64 {
	return rand.Int63n(4294967295)
}

func SendPWResetCode(to_email string, code int64) {
	msg := mail.NewMessage()
	msg.SetHeader("From", config.Vals.Server.SMTPUser)
	msg.SetHeader("To", to_email)
	msg.SetHeader("Subject", "Comment Anywhere - Password Reset Code")
	msg.SetBody("text/html", "<h2>Your password reset code is: <br /><b>"+fmt.Sprintf("%d", code)+"</b></h2>")

	d := mail.NewDialer(config.Vals.Server.SMTPServer, 587, config.Vals.Server.SMTPUser, config.Vals.Server.SMTPPass)

	d.StartTLSPolicy = mail.MandatoryStartTLS

	err := d.DialAndSend(msg)
	if err != nil {
		fmt.Printf("\nError; failed to send pw reset email to: %s", to_email)
	}
}

func SendVerificationCode(to_email string, username string, code int64) {
	msg := mail.NewMessage()
	msg.SetHeader("From", config.Vals.Server.SMTPUser)
	msg.SetHeader("To", to_email)
	msg.SetHeader("Subject", "Comment Anywhere - Verification Code")
	msg.SetBody("text/html", "Welcome to comment anywhere, "+username+". \n\nYour account verification code is: <b>"+fmt.Sprintf("%d", code)+"</b>")

	d := mail.NewDialer(config.Vals.Server.SMTPServer, 587, config.Vals.Server.SMTPUser, config.Vals.Server.SMTPPass)

	d.StartTLSPolicy = mail.MandatoryStartTLS

	err := d.DialAndSend(msg)
	if err != nil {
		fmt.Printf("\nError; failed to send account verification email to: %s", to_email)
	}

}
