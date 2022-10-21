package feedback

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/justindfuller/interviews/organization"
)

func GiveHandler(organizations *organization.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			cookie, err := r.Cookie("__Host-UserUUID")
			if err != nil {
				log.Printf("Error parsing cookie for /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/unauthenticated.html")
				return
			}

			if cookie.Value == "" {
				log.Printf("Error parsing cookie for /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/unauthenticated.html")
				return
			}

			org, err := organizations.FindByUserID(cookie.Value)
			if err != nil {
				log.Printf("Error finding organization for /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			paths := strings.Split(r.URL.Path, "/")

			if len(paths) < 4 || paths[3] == "" {
				log.Printf("No ID provided for /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			id, err := uuid.Parse(paths[3])
			if err != nil {
				log.Printf("Error parsing ID for /feedback/give/{id}: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			f, err := org.FeedbackByID(id)
			if err != nil {
				log.Printf("Error finding feedback for  /feedback/give/%s: %s", id, err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			t, err := template.ParseFiles("./feedback/give.html")
			if err != nil {
				log.Printf("Error parsing template for /feedback/give/%s: %s", id, err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := t.Execute(w, f); err != nil {
				log.Printf("Error executing template for /feedback/give/%s: %s", id, err)
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

			log.Printf("Received /feedback/give: %#v", query)
			return
		}

		log.Printf("Unexpected http Method '%s' for /feedback/give", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
