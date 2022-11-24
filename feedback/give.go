package feedback

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	interview "github.com/justindfuller/interviews"
)

const GivePath = "/feedback/give/"

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

		f, request, err := org.FeedbackByRequestID(id)
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
			t, err := template.New("give.template.html").Funcs(funcs).ParseFiles("./feedback/give.template.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /feedback/give/%s: %s", id, err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			vars := map[string]interface{}{
				"FeedbackRequestID": request.ID,
				"Candidate":         request.CandidateName,
				"CreatorID":         f.CreatorID,
				"CreatedAt":         f.CreatedAt,
				"Role":              f.Role,
				"Team":              f.Team,
				"Questions":         f.Questions,
			}
			if err := t.Execute(w, vars); err != nil {
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
				if key == "recommend" {
					continue
				}

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

			recommend, err := strconv.ParseBool(query.Get("recommend"))
			if err != nil {
				log.Printf("Error parsing recommend boolean in /feedback/give: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			given, err := interview.NewFeedbackResponse(userID, answers, recommend)
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

			go func() {
				funcs := template.FuncMap{
					"UserEmail": func(id uuid.UUID) (string, error) {
						user, err := org.FindUserByID(id.String())
						if err != nil {
							return "", err
						}
						return user.Email, nil
					},
				}
				t, err := template.New("given-email.template.html").Funcs(funcs).ParseFiles("./feedback/given-email.template.html", "index.css")
				if err != nil {
					log.Printf("Error parsing invite template for /feedback/: %s", err)
					return
				}

				var html strings.Builder
				vars := map[string]interface{}{
					"GiverID":   userID,
					"Team":      f.Team,
					"Role":      f.Role,
					"Questions": f.Questions,
					"Date":      time.Now(),
					"Responses": []interview.FeedbackResponse{given},
					"Recommend": given.Recommend,
				}
				if err := t.Execute(&html, vars); err != nil {
					log.Printf("Error executing invite template for /feedback/: %s", err)
				}

				user, err := org.FindUserByID(f.CreatorID.String())
				if err != nil {
					log.Printf("Error finding user for /feedback/: %s", err)
					return
				}

				opts := interview.EmailOptions{
					To:      []string{user.Email},
					Subject: fmt.Sprintf("Feedback received for the role %s on team %s", f.Role, f.Team),
					HTML:    html.String(),
				}
				if err := interview.Email(opts); err != nil {
					log.Printf("Error sending email from /feedback/give: %s", err)
				}
			}()

			http.Redirect(w, r, "/organization/", http.StatusSeeOther)
			return
		}

		log.Printf("Unexpected http Method '%s' for /feedback/give", r.Method)
		http.ServeFile(w, r, "./error/index.html")
	}
}
