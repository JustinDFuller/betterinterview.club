package organization

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

const InvitePath = "/organization/invite/"

func InviteHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil {
			log.Printf("Error parsing cookie for /organization/invite: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if cookie.Value == "" {
			log.Printf("Error parsing cookie for /organization/invite: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		org, inviter, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /organization/invite: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if r.Method == http.MethodGet {
			t, err := template.New("invite.template.html").ParseFiles("./organization/invite.template.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			vars := map[string]interface{}{
				"Domain": org.Domain,
				"Error":  "",
				"Email":  "",
			}
			if err := t.Execute(w, vars); err != nil {
				log.Printf("Error executing template for /organization: %s", err)
				return
			}

			return
		}

		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading /organization/invite body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}
			defer r.Body.Close()

			query, err := url.ParseQuery(string(body))
			if err != nil {
				log.Printf("Error parsing query from /organization/invite body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			email, err := mail.ParseAddress(query.Get("email"))
			if err != nil {
				log.Printf("Error parsing email from /organization/invite email parameter: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			isDifferentDomain, err := org.IsDifferentDomain(email.Address)
			if err != nil {
				log.Printf("Error validating email domain in /organization/invite: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}
			if isDifferentDomain {
				t, err := template.New("invite.template.html").ParseFiles("./organization/invite.template.html", "index.css")
				if err != nil {
					log.Printf("Error parsing template for /organization: %s", err)
					http.ServeFile(w, r, "./error/index.html")
					return
				}

				vars := map[string]interface{}{
					"Domain": org.Domain,
					"Error":  "You cannot invite members from another domain.",
					"Email":  email.Address,
				}
				if err := t.Execute(w, vars); err != nil {
					log.Printf("Error executing template for /organization: %s", err)
					return
				}

				return
			}

			org, user, err := organizations.FindOrCreateByEmail(email.Address)
			if err != nil {
				log.Printf("Error finding or creating org for user in /organization/invite: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			cbID, err := organizations.AddEmailLoginCallback(org, user)
			if err != nil {
				log.Printf("Error adding email login callback in /organization/invite: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			go func() {
				t, err := template.New("invite-email.template.html").ParseFiles("./organization/invite-email.template.html", "index.css")
				if err != nil {
					log.Printf("Error parsing invite template for /organization/invite/: %s", err)
					return
				}

				var html strings.Builder
				vars := map[string]string{
					"ID":      cbID,
					"Host":    os.Getenv("HOST"),
					"Inviter": inviter.Email,
				}
				if err := t.Execute(&html, vars); err != nil {
					log.Printf("Error executing invite template for /organization/invite/: %s", err)
				}

				opts := interview.EmailOptions{
					To:      email.Address,
					Subject: "Your invite",
					HTML:    html.String(),
				}
				if err := interview.Email(opts, org); err != nil {
					log.Printf("Error sending email from /organization/invite: %s", err)
					return
				}
			}()
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}
}
