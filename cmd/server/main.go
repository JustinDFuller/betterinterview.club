package main

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/feedback"
	"github.com/justindfuller/interviews/organization"
)

func main() {
	var organizations organization.Organizations

	http.HandleFunc("/feedback/", feedback.Handler(&organizations))
	http.HandleFunc("/organization/member/", organization.MemberHandler(&organizations))
	http.HandleFunc("/organization/", organization.Handler(&organizations))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("__Host-UserUUID")
		if err != nil || cookie == nil || cookie.Value == "" {
			http.ServeFile(w, r, "./index.html")
			return
		}

		if _, err := organizations.FindByUserID(cookie.Value); err != nil {
			log.Printf("Error finding organization for /: %s", err)
			http.ServeFile(w, r, "./index.html")
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	})

	log.Print("Listening at http://localhost:8443/")
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))
}
