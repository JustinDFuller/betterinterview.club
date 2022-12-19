package auth

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"strings"

	interview "github.com/justindfuller/interviews"
)

const LoginPath = "/auth/login/"

var commonEmails = []string{
	"gmail.com",
	"yahoo.com",
	"ymail.com",
	"hotmail.com",
	"msn.com",
	"live.com",
	"outlook.com",
	"verizon.net",
	"icloud.com",
	"att.net",
	"mac.com",
}

var (
	loginTemplate      = template.Must(template.New("login.template.html").ParseFiles("./auth/login.template.html", "index.css"))
	loginEmailTemplate = template.Must(template.New("login-email.template.html").ParseFiles("./auth/login-email.template.html", "index.css"))
)

func LoginHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := loginTemplate.Execute(w, nil); err != nil {
				log.Printf("Error executing template for /auth/login: %s", err)
			}

			return
		}

		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading /auth/login body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}
			defer r.Body.Close()

			query, err := url.ParseQuery(string(body))
			if err != nil {
				log.Printf("Error parsing query from /auth/login body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			email, err := mail.ParseAddress(query.Get("email"))
			if err != nil {
				log.Printf("Error parsing email from /auth/login email parameter: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			split := strings.Split(email.Address, "@")
			if len(split) != 2 {
				log.Printf("Invalid email address: incorrect number of parts: %s", email.Address)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			domain := strings.ToLower(split[1])

			for _, common := range commonEmails {
				if domain == common {
					vars := map[string]interface{}{
						"Error": "Public email domains are not allowed.",
					}
					if err := loginTemplate.Execute(w, vars); err != nil {
						log.Printf("Error executing template for /auth/login: %s", err)
					}
					return
				}
			}

			org, user, err := organizations.FindOrCreateByEmail(email.Address)
			if err != nil {
				log.Printf("Error finding or creating org for user in /auth/login: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			cbID, err := organizations.AddEmailLoginCallback(org, user)
			if err != nil {
				log.Printf("Error adding email login callback in /auth/login: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			var html strings.Builder
			if err := loginEmailTemplate.Execute(&html, map[string]string{"ID": cbID, "Host": os.Getenv("HOST")}); err != nil {
				log.Printf("Error executing template for /auth/login: %s", err)
			}

			go func() {
				opts := interview.EmailOptions{
					To:      email.Address,
					Subject: "Log in",
					HTML:    html.String(),
				}
				if err := interview.Email(opts, org); err != nil {
					log.Printf("Error sending email from /auth/login: %s", err)
				}
			}()

			http.Redirect(w, r, "/auth/email/", http.StatusSeeOther)
			return
		}

		log.Printf("Unexpected http Method '%s' for /auth/login", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
