package api

import (
	"log"
	"net/http"
	"text/template"

	interview "github.com/justindfuller/interviews"
	"github.com/justindfuller/interviews/auth"
	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

func Handlers() {
	var organizations interview.Organizations

	http.HandleFunc("/auth/login/", auth.LoginHandler(&organizations))
	http.HandleFunc("/auth/callback/", auth.CallbackHandler(&organizations))
	http.HandleFunc("/auth/logout/", auth.LogoutHandler)
	http.HandleFunc("/auth/email/", auth.EmailHandler)
	http.HandleFunc("/feedback/given/", feedback.GivenHandler(&organizations))
	http.HandleFunc("/feedback/give/", feedback.GiveHandler(&organizations))
	http.HandleFunc("/feedback/", feedback.Handler(&organizations))
	// http.HandleFunc("/organization/notfound/", organization.NotFoundHandler)
	http.HandleFunc("/organization/member/", organization.MemberHandler(&organizations))
	http.HandleFunc("/organization/", organization.Handler(&organizations))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			t, err := template.New("login.html").ParseFiles("auth/login.html", "index.css")
			if err != nil {
				log.Printf("Error parsing template for /: %s", err)
				http.ServeFile(w, r, "./error/index.html")
				return
			}

			if err := t.Execute(w, nil); err != nil {
				log.Printf("Error executing template for /: %s", err)
			}
			return
		}

		if _, err := organizations.FindByUserID(cookie.Value); err != nil {
			auth.LoginHandler(&organizations)(w, r)
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	})
}
