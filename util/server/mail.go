package server

import "net/smtp"

/* See: https://rfc-editor.org/rfc/rfc5321.html */

// Dispatches an email.
func SendMail(address string, from string, subject string, body string, to string) error {
	client, err := smtp.Dial("localhost:25")
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Mail(from)
	if err != nil {
		return err
	}
	err = client.Rcpt(to)
	if err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(body))
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	client.Quit()
	return nil
}
