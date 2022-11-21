package api

import (
	"compress/gzip"
	"html/template"
	"io"
	"log"
	"net/http"

	interview "github.com/justindfuller/interviews"
	"github.com/justindfuller/interviews/auth"
	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

type gzipResponseWriter struct {
	w      http.ResponseWriter
	writer io.Writer
}

func (grw gzipResponseWriter) Header() http.Header {
	return grw.w.Header()
}

func (grw gzipResponseWriter) WriteHeader(statusCode int) {
	grw.w.WriteHeader(statusCode)
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.writer.Write(b)
}

func withGzip(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := gzip.NewWriter(w)
		defer writer.Close()

		gzr := gzipResponseWriter{w, writer}
		w.Header().Set("Content-Encoding", "gzip")
		handler(gzr, r)
	}
}

func Handlers() {
	var organizations interview.Organizations

	http.HandleFunc("/auth/login/", withGzip(auth.LoginHandler(&organizations)))
	http.HandleFunc("/auth/callback/", withGzip(auth.CallbackHandler(&organizations)))
	http.HandleFunc("/auth/logout/", withGzip(auth.LogoutHandler))
	http.HandleFunc("/auth/email/", withGzip(auth.EmailHandler(&organizations)))
	http.HandleFunc("/feedback/given/", withGzip(feedback.GivenHandler(&organizations)))
	http.HandleFunc("/feedback/give/", withGzip(feedback.GiveHandler(&organizations)))
	http.HandleFunc("/feedback/close/", withGzip(feedback.CloseHandler(&organizations)))
	http.HandleFunc("/feedback/", withGzip(feedback.Handler(&organizations)))
	http.HandleFunc("/organization/member/", withGzip(organization.MemberHandler(&organizations)))
	http.HandleFunc("/organization/", withGzip(organization.Handler(&organizations)))

	http.HandleFunc("/", withGzip(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			t, err := template.New("login.template.html").ParseFiles("auth/login.template.html", "index.css")
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

		if _, _, err := organizations.FindByUserID(cookie.Value); err != nil {
			auth.LoginHandler(&organizations)(w, r)
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}))
}
