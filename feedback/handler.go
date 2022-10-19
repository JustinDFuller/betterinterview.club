package feedback

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/organization"
)

func Handler(organizations *organization.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil {
			log.Printf("Error parsing cookie for /feedback: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if cookie.Value == "" {
			log.Printf("Error parsing cookie for /feedback: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if _, err := organizations.FindByUserID(cookie.Value); err != nil {
			log.Printf("Error finding organization for /feedback: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		http.ServeFile(w, r, "./feedback/index.html")
	}
}
