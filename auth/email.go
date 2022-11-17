package auth

import (
	"html/template"
	"log"
	"net/http"

	interview "github.com/justindfuller/interviews"
)

func EmailHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err == nil {
			if _, _, err := organizations.FindByUserID(cookie.Value); err == nil {
				http.Redirect(w, r, "/organization/", http.StatusSeeOther)
				return
			}
		}

		t, err := template.New("email.html").ParseFiles("./auth/email.html", "index.css")
		if err != nil {
			log.Printf("Error parsing template for /auth/email/: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Printf("Error executing template for /auth/email: %s", err)
		}
	}
}
