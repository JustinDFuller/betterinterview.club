package api

import (
	"net/http"

	interview "github.com/justindfuller/interviews"
	"github.com/justindfuller/interviews/auth"
	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

func Handlers() {
	var organizations interview.Organizations

	http.HandleFunc(auth.LoginPath, withGzip(auth.LoginHandler(&organizations)))
	http.HandleFunc(auth.CallbackPath, withGzip(auth.CallbackHandler(&organizations)))
	http.HandleFunc(auth.LogoutPath, withGzip(auth.LogoutHandler))
	http.HandleFunc(auth.EmailPath, withGzip(auth.EmailHandler(&organizations)))
	http.HandleFunc(feedback.GivenPath, withGzip(feedback.GivenHandler(&organizations)))
	http.HandleFunc(feedback.GivePath, withGzip(feedback.GiveHandler(&organizations)))
	http.HandleFunc(feedback.ClosePath, withGzip(feedback.CloseHandler(&organizations)))
	http.HandleFunc(feedback.OpenPath, withGzip(feedback.OpenHandler(&organizations)))
	http.HandleFunc(organization.InvitePath, withGzip(organization.InviteHandler(&organizations)))
	http.HandleFunc(organization.Path, withGzip(organization.Handler(&organizations)))

	http.HandleFunc("/", withGzip(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			auth.LoginHandler(&organizations)(w, r)
			return
		}

		if _, _, err := organizations.FindByUserID(cookie.Value); err != nil {
			auth.LoginHandler(&organizations)(w, r)
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}))
}
