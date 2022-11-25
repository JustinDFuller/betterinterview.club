package feedback

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

const OpenPath = "/feedback/"

func OpenHandler(organizations *interview.Organizations) http.HandlerFunc {
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

		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			log.Printf("Error parsing cookie for /feedback: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		org, _, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /feedback: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if r.Method == http.MethodGet {
			t, err := template.New("open.template.html").ParseFiles("feedback/open.template.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /: %s", err)
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

			q1, err := interview.NewQuestion(query.Get("question1"))
			if err != nil {
				log.Printf("Error creating question from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			q2, err := interview.NewQuestion(query.Get("question2"))
			if err != nil {
				log.Printf("Error creating question from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			q3, err := interview.NewQuestion(query.Get("question3"))
			if err != nil {
				log.Printf("Error creating question from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			q4, err := interview.NewQuestion(query.Get("question4"))
			if err != nil {
				log.Printf("Error creating question from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			q5, err := interview.NewQuestion(query.Get("question5"))
			if err != nil {
				log.Printf("Error creating question from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			f, err := interview.NewFeedback(userID, query.Get("team"), query.Get("role"), []interview.Question{q1, q2, q3, q4, q5})
			if err != nil {
				log.Printf("Error creating feedback from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := organizations.AddFeedback(org, f); err != nil {
				log.Printf("Error adding feedback to organization")
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			log.Printf("New Feedback: %s %s", org.Domain, f)
			http.Redirect(w, r, "/organization/", http.StatusSeeOther)
			return
		}

		log.Printf("Unexpected http Method '%s' for /feedback", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
