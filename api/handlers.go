package api

import (
	"net/http"

	interview "github.com/justindfuller/interviews"
	"github.com/justindfuller/interviews/auth"
	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

func Handlers() {
	organizations := interview.DefaultOrganizations

	http.HandleFunc(auth.LoginPath, middleware(auth.LoginHandler(organizations)))
	http.HandleFunc(auth.CallbackPath, middleware(auth.CallbackHandler(organizations)))
	http.HandleFunc(auth.LogoutPath, middleware(auth.LogoutHandler))
	http.HandleFunc(auth.EmailPath, middleware(auth.EmailHandler(organizations)))
	http.HandleFunc(feedback.RequestPath, middleware(feedback.RequestHandler(organizations)))
	http.HandleFunc(feedback.GivenPath, middleware(feedback.GivenHandler(organizations)))
	http.HandleFunc(feedback.GivePath, middleware(feedback.GiveHandler(organizations)))
	http.HandleFunc(feedback.ClosePath, middleware(feedback.CloseHandler(organizations)))
	http.HandleFunc(feedback.OpenPath, middleware(feedback.OpenHandler(organizations)))
	http.HandleFunc(organization.InvitePath, middleware(organization.InviteHandler(organizations)))
	http.HandleFunc(organization.Path, middleware(organization.Handler(organizations)))

	http.HandleFunc("/", middleware(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			auth.LoginHandler(organizations)(w, r)
			return
		}

		if _, _, err := organizations.FindByUserID(cookie.Value); err != nil {
			auth.LoginHandler(organizations)(w, r)
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}))
}
