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

		t, err := template.New("invite.html").ParseFiles("./organization/invite.html", "index.css")
		if err != nil {
			log.Printf("Error parsing invite template for /organization/member/: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		var html strings.Builder
		if err := t.Execute(&html, map[string]string{"ID": cbID, "Host": os.Getenv("HOST")}); err != nil {
			log.Printf("Error executing invite template for /organization/member/: %s", err)
		}

		opts := interview.EmailOptions{
			To:      email.Address,
			Subject: "Your invite to Better Interviews",
			HTML:    html.String(),
		}
		if err := interview.Email(opts); err != nil {
			log.Printf("Error sending email from /auth/login: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}
}
