package auth

import (
	"html/template"
	"log"
	"net/http"

	interview "github.com/justindfuller/interviews"
)

const EmailPath = "/auth/email/"

var emailTemplate = template.Must(template.New("email.template.html").ParseFiles("./auth/email.template.html", "index.css"))

func EmailHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err == nil {
			if _, _, err := organizations.FindByUserID(cookie.Value); err == nil {
				http.Redirect(w, r, "/organization/", http.StatusSeeOther)
				return
			}
		}

		if err := emailTemplate.Execute(w, nil); err != nil {
			log.Printf("Error executing template for /auth/email: %s", err)
		}
	}
}
