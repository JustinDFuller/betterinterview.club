package feedback

import (
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

const RequestPath = "/feedback/request/"

func RequestHandler(organizations *interview.Organizations) http.HandlerFunc {
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

		org, creator, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /feedback: %s", err)
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

		if r.Method == http.MethodGet {
			t, err := template.New("request.template.html").ParseFiles("feedback/request.template.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			vars := map[string]interface{}{
				"Feedback": f,
				"Domain":   org.Domain,
			}
			if err := t.Execute(w, vars); err != nil {
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

			request, err := interview.NewFeedbackRequest(query.Get("candidate"), query.Get("email1"), query.Get("email2"))
			if err != nil {
				log.Printf("Error creating feedback request: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := organizations.AddFeedbackRequest(org, f, request); err != nil {
				log.Printf("Error adding feedback to organization: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			for _, email := range request.InterviewerEmails {
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
						"ID":                cbID,
						"FeedbackRequestID": request.ID.String(),
						"Host":              os.Getenv("HOST"),
						"Role":              f.Role,
						"Team":              f.Team,
						"Candidate":         request.CandidateName,
						"RequestedBy":       creator.Email,
					}
					if err := t.Execute(&html, vars); err != nil {
						log.Printf("Error executing invite template for /feedback/: %s", err)
					}

					opts := interview.EmailOptions{
						To:      []string{email},
						Subject: "Feedback Requested",
						HTML:    html.String(),
					}
					if err := interview.Email(opts); err != nil {
						log.Printf("Error sending email from /auth/login: %s", err)
					}
				}()
			}

			log.Printf("New Feedback: %s %s", org.Domain, f)
			http.Redirect(w, r, "/organization/", http.StatusSeeOther)

			return

		}

		log.Printf("Unexpected http Method '%s' for /feedback", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
