package interview

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/pkg/errors"
)

type EmailOptions struct {
	To      string
	Subject string
	HTML    string
}

func Email(opts EmailOptions) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte(fmt.Sprintf("To: %s\r\n", opts.To) +
		fmt.Sprintf("Subject: %s\r\n", opts.Subject) +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		"<!DOCTYPE html><html><body>\r\n" +
		opts.HTML +
		"</body></html>")
	err := smtp.SendMail("smtp.gmail.com:587", auth, email, []string{opts.To}, msg)
	if err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}
