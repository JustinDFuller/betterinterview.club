package auth

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/google/uuid"
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

			cbID, err := uuid.NewRandom()
			if err != nil {
				log.Printf("Error creating callback ID in /auth/login: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			opts := interview.EmailOptions{
				To:      email.Address,
				Subject: "Log in to Better Interviews",
				HTML:    fmt.Sprintf(`<h1>Better Interviews</h1><a href="https://localhost:8443/auth/callback?id=%s">Log In</a>`, cbID.String()),
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

/*
org, err := organizations.FindByUserEmail(email.Address)
			if err != nil {
				log.Printf("Unable to find organization /auth/login: %s", err)
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, "./organization/notfound.html")
				return
			}

			user, err := org.FindUserByEmail(email.Address)
			if err != nil {
				log.Printf("Unable to find user /auth/login: %s", err)
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, "./organization/notfound.html")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "__Host-UserUUID",
				Value:    user.ID.String(),
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24 * 31), // One month
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			http.Redirect(w, r, "/organization/", http.StatusSeeOther)
*/
