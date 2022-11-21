package feedback

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

func CloseHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil {
			log.Printf("Error parsing cookie for /feedback/close: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if cookie.Value == "" {
			log.Printf("Error parsing cookie for /feedback/close: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		/*userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			log.Printf("Error parsing cookie for /feedback/close: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}*/

		org, closer, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /feedback/close: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		paths := strings.Split(r.URL.Path, "/")

		if len(paths) < 4 || paths[3] == "" {
			log.Printf("No ID provided for /feedback/close: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		id, err := uuid.Parse(paths[3])
		if err != nil {
			log.Printf("Error parsing ID for /feedback/close/{id}: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		f, err := org.FeedbackByID(id)
		if err != nil {
			log.Printf("Error finding feedback for  /feedback/close/%s: %s", id, err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if f.CreatorID != closer.ID {
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		if r.Method == http.MethodGet {
			t, err := template.New("close.template.html").ParseFiles("./feedback/close.template.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /feedback/close/%s: %s", id, err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			vars := map[string]string{
				"ID":        f.ID.String(),
				"Role":      f.Role,
				"Team":      f.Team,
				"Email":     closer.Email,
				"CreatedAt": f.CreatedAt.Format("January 2 2006"),
			}
			if err := t.Execute(w, vars); err != nil {
				log.Printf("Error executing template for /feedback/close/%s: %s", id, err)
			}

			return
		}

		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading /feedback/give body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}
			defer r.Body.Close()

			query, err := url.ParseQuery(string(body))
			if err != nil {
				log.Printf("Error parsing query from /feedback/give body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			f.Closed = true
			f.CloseReason = query.Get("reason")

			if err := organizations.SetFeedback(org, f); err != nil {
				log.Printf("Error saving feedback /feedback/close/%s: %s", id, err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			http.Redirect(w, r, "/organization/", http.StatusSeeOther)
		}

		log.Printf("Unexpected http Method '%s' for /feedback/close", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
