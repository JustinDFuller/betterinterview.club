package organization

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

const Path = "/organization/"

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

		org, user, err := organizations.FindByUserID(cookie.Value)
		if err != nil {
			log.Printf("Error finding organization for /organization: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		funcs := template.FuncMap{
			"UserEmail": func(id uuid.UUID) interview.User {
				user, _ := org.FindUserByID(id.String())
				return user
			},
		}
		t, err := template.New("index.template.html").Funcs(funcs).ParseFiles("./organization/index.template.html", "index.css")
		if err != nil {
			log.Printf("Error parsing template for /organization: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		var feedback []interview.Feedback
		for _, f := range org.Feedback {
			if !f.Closed && f.CreatorID == user.ID {
				feedback = append(feedback, f)
			}
		}

		var requests []interview.FeedbackRequest
		for _, f := range org.Feedback {
			if f.Closed {
				continue
			}

			for _, request := range f.Requests {
				var emailMatch bool

				for _, email := range request.InterviewerEmails {
					if email == user.Email {
						emailMatch = true
					}
				}

				if !emailMatch {
					continue
				}

				var requestComplete bool

				for _, response := range request.Responses {
					if response.CreatorID == user.ID {
						requestComplete = true
					}
				}

				if requestComplete {
					continue
				}

				requests = append(requests, request)
			}
		}

		vars := map[string]interface{}{
			"UserID":   user.ID.String(),
			"Feedback": feedback,
			"Requests": requests,
			"Users":    org.Users,
			"Domain":   org.Domain,
		}
		if err := t.Execute(w, vars); err != nil {
			log.Printf("Error executing template for /organization: %s", err)
		}

		return
	}
}
