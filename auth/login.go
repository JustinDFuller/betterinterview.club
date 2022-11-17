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

func LoginHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			t, err := template.New("login.html").ParseFiles("./auth/login.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /auth/login/: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := t.Execute(w, nil); err != nil {
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

			t, err := template.New("login-email.html").ParseFiles("./auth/login-email.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /auth/login/: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			var html strings.Builder
			if err := t.Execute(&html, map[string]string{"ID": cbID, "Host": os.Getenv("HOST")}); err != nil {
				log.Printf("Error executing template for /auth/login: %s", err)
			}

			opts := interview.EmailOptions{
				To:      []string{email.Address},
				Subject: "Log in to Better Interviews",
				HTML:    html.String(),
			}
			if err := interview.Email(opts); err != nil {
				log.Printf("Error sending email from /auth/login: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			http.Redirect(w, r, "/auth/email/", http.StatusSeeOther)
			return
		}

		log.Printf("Unexpected http Method '%s' for /auth/login", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
