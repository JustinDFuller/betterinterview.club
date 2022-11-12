package interview

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/pkg/errors"
)

func Email(to, subject, html string) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" +
		"<!DOCTYPE html><html><body>\r\n" +
		html +
		"</body></html>")
	err := smtp.SendMail("smtp.gmail.com:587", auth, email, []string{to}, msg)
	if err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}
