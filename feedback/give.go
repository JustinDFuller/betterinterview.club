package feedback

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

func GiveHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			log.Printf("Error parsing cookie for /feedback/give: %s", err)
			http.ServeFile(w, r, "./error/unauthenticated.html")
			return
		}

		org, _, err := organizations.FindByUserID(cookie.Value)
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

		if r.Method == http.MethodGet {
			funcs := template.FuncMap{
				"UserEmail": func(id uuid.UUID) (string, error) {
					user, err := org.FindUserByID(id.String())
					if err != nil {
						return "", err
					}
					return user.Email, nil
				},
			}
			t, err := template.New("give.html").Funcs(funcs).ParseFiles("./feedback/give.html", "index.css")
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

			var answers []interview.Answer

			for key := range query {
				b, err := strconv.ParseBool(query.Get(key))
				if err != nil {
					log.Printf("Error parsing Answer boolean in /feedback/give: %s", err)
					http.ServeFile(w, r, "./error/index.html")
					return
				}

				a, err := interview.NewAnswer(key, b)
				if err != nil {
					log.Printf("Error creating Answer in /feedback/give: %s", err)
					http.ServeFile(w, r, "./error/index.html")
					return
				}

				answers = append(answers, a)
			}

			given, err := interview.NewFeedbackResponse(userID, answers)
			if err != nil {
				log.Printf("Error creating FeedbackResponse in /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := organizations.AddFeedbackResponse(org, f, given); err != nil {
				log.Printf("Error adding feedback response in /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			org, _, err = organizations.FindByUserID(cookie.Value)
			if err != nil {
				log.Printf("Error finding organization for /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			http.Redirect(w, r, "/organization/", http.StatusSeeOther)
			return
		}

		log.Printf("Unexpected http Method '%s' for /feedback/give", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
