package interview

import (
	"bytes"
	"net/smtp"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

type EmailOptions struct {
	To      string
	Subject string
	HTML    string
	From    string
}

func Email(opts EmailOptions) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	opts.From = email

	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	t, err := template.New("email.template.txt").ParseFiles("./email.template.txt")
	if err != nil {
		return errors.Wrap(err, "error parsing emplate.template")
	}

	var b bytes.Buffer
	if err := t.Execute(&b, opts); err != nil {
		return errors.Wrap(err, "error executing email.template")
	}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, email, []string{opts.To}, b.Bytes()); err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}
