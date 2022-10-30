package organization

import (
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

func MemberHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil {
			log.Printf("Error parsing cookie for /organization/member: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if cookie.Value == "" {
			log.Printf("Error parsing cookie for /organization/member: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		org, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /organization/member: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading /organization/member body: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}
		defer r.Body.Close()

		query, err := url.ParseQuery(string(body))
		if err != nil {
			log.Printf("Error parsing query from /organization/member body: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		email, err := mail.ParseAddress(query.Get("email"))
		if err != nil {
			log.Printf("Error parsing email from /organization/member email parameter: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		parts := strings.Split(email.Address, "@")
		if len(parts) < 2 {
			log.Printf("Error parsing email from /organization/member email parameter: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if parts[1] != org.Domain {
			// TODO: This is needs a better error for the user.
			log.Printf("Invalid email in /organization/member email parameter: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		userID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("Error creating user UUID in /organization/member: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		u := interview.User{
			ID:    userID,
			Email: email.Address,
		}
		if _, err = organizations.AddUser(org, u); err != nil {
			log.Printf("Error adding user to organization in /organization/member: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}
}
