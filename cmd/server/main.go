package main

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

func main() {
	http.HandleFunc("/feedback/", feedback.Handler)
	http.HandleFunc("/organization/member/", organization.MemberHandler)
	http.HandleFunc("/organization/", organization.Handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			http.ServeFile(w, r, "./index.html")
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	})

	log.Print("Listening at http://localhost:8443/")
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))
}
