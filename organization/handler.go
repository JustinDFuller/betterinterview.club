package organization

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

func Handler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cookie, err := r.Cookie("__Host-UserUUID")
			if err != nil {
				log.Printf("Error parsing cookie for /organization: %s", err)
				http.ServeFile(w, r, "./error/unauthenticated.html")
				return
			}

			if cookie.Value == "" {
				log.Printf("Error parsing cookie for /organization: %s", err)
				http.ServeFile(w, r, "./error/unauthenticated.html")
				return
			}

			org, err := organizations.FindByUserID(cookie.Value)
			if err != nil {
				log.Printf("Error finding organization for /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			funcs := template.FuncMap{
				"UserEmail": func(id uuid.UUID) (interview.User, error) {
					return org.FindUserByID(id.String())
				},
			}
			t, err := template.New("index.html").Funcs(funcs).ParseFiles("./organization/index.html")
			if err != nil {
				log.Printf("Error parsing template for /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := t.Execute(w, org); err != nil {
				log.Printf("Error executing template for /organization: %s", err)
			}

			return
		}

		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading /organization body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}
			defer r.Body.Close()

			query, err := url.ParseQuery(string(body))
			if err != nil {
				log.Printf("Error parsing query from /organization body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			email, err := mail.ParseAddress(query.Get("email"))
			if err != nil {
				log.Printf("Error parsing email from /organization email parameter: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			parts := strings.Split(email.Address, "@")
			if len(parts) < 2 {
				log.Printf("Error parsing email from /organization email parameter: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			orgID, err := uuid.NewRandom()
			if err != nil {
				log.Printf("Error creating org UUID in /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			user, err := interview.NewUser(email.Address)
			if err != nil {
				log.Printf("Error creating user UUID in /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			org, err := organizations.FindByDomain(email.Address)
			if err != nil {
				org = interview.Organization{
					ID:     orgID,
					Domain: parts[1],
				}
				if err := organizations.Add(org); err != nil {
					log.Printf("Error adding organization '%s' from /organization: %s", org.Domain, err)
					http.ServeFile(w, r, "./organization/exists.html")
					return
				}
			}

			if _, err := organizations.AddUser(org, user); err != nil {
				log.Printf("Error adding user in /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
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

			log.Printf("New Organization: %s %#v", org.Domain, org)

			http.Redirect(w, r, "/organization/", http.StatusSeeOther)

			return
		}

		log.Printf("Unexpected http Method '%s' for /organization", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
