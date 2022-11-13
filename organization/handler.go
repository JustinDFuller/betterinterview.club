package organization

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

func Handler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		t, err := template.New("index.html").Funcs(funcs).ParseFiles("./organization/index.html", "index.css")
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
}
