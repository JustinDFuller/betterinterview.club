package interview

import (
	"bytes"
	"log"
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

var ErrCrossDomainEmail = errors.New("cannot send cross-domain emails")

var emailTemplate = template.Must(template.New("email.template.txt").ParseFiles("./email.template.txt"))

func Email(opts EmailOptions, org Organization) error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	isDifferentDomain, err := org.IsDifferentDomain(opts.To)
	log.Printf("Domain: %s, To: %s, IsDifferent: %v, Err: %s", org.Domain, opts.To, isDifferentDomain, err)
	if err != nil {
		return errors.Wrap(err, "error validating email domain")
	}
	if isDifferentDomain {
		return ErrCrossDomainEmail
	}

	// Adding to opts so that it appears in the template.
	opts.From = email

	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	var b bytes.Buffer
	if err := emailTemplate.Execute(&b, opts); err != nil {
		return errors.Wrap(err, "error executing email.template")
	}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, email, []string{opts.To}, b.Bytes()); err != nil {
		return errors.Wrap(err, "error sending email")
	}

	return nil
}
