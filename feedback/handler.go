package feedback

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

		org, creator, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /feedback: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		if r.Method == http.MethodGet {
			t, err := template.New("index.template.html").ParseFiles("feedback/index.template.html", "index.css")
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

			var emails []string
			if email := query.Get("email1"); email != "" {
				emails = append(emails, email)
			}
			if email := query.Get("email2"); email != "" {
				emails = append(emails, email)
			}

			f, err := interview.NewFeedback(userID, query.Get("team"), query.Get("role"), emails, []interview.Question{q1, q2, q3, q4, q5})
			if err != nil {
				log.Printf("Error creating feedback from /feedback body: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			for _, email := range emails {
				email := email

				go func() {
					org, user, err := organizations.FindOrCreateByEmail(email)
					if err != nil {
						log.Printf("Error finding or creating org for invited user in /feedback/: %s", err)
						http.ServeFile(w, r, "./error/index.html")
						return
					}

					cbID, err := organizations.AddEmailLoginCallback(org, user)
					if err != nil {
						log.Printf("Error adding email login callback in /feedback/: %s", err)
						http.ServeFile(w, r, "./error/index.html")
						return
					}

					t, err := template.New("invite.template.html").ParseFiles("./feedback/invite.template.html", "index.css")
					if err != nil {
						log.Printf("Error parsing invite template for /feedback/: %s", err)
						return
					}

					var html strings.Builder
					vars := map[string]string{
						"ID":          cbID,
						"FeedbackID":  f.ID.String(),
						"Host":        os.Getenv("HOST"),
						"Role":        query.Get("role"),
						"Team":        query.Get("team"),
						"RequestedBy": creator.Email,
					}
					if err := t.Execute(&html, vars); err != nil {
						log.Printf("Error executing invite template for /feedback/: %s", err)
					}

					opts := interview.EmailOptions{
						To:      []string{email},
						Subject: fmt.Sprintf("Your feedback is requested for the role %s on team %s", query.Get("role"), query.Get("team")),
						HTML:    html.String(),
					}
					if err := interview.Email(opts); err != nil {
						log.Printf("Error sending email from /auth/login: %s", err)
					}
				}()
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
