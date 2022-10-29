package auth

import (
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/justindfuller/interviews/organization"
)

func LoginHandler(organizations *organization.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "./auth/login.html")
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

			org, err := organizations.FindByUserEmail(email.Address)
			if err != nil {
				log.Printf("Unable to find organization /auth/login: %s", err)
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, "./organization/notfound.html")
				return
			}

			log.Printf("%s", org)
			return
		}

		log.Printf("Unexpected http Method '%s' for /auth/login", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
