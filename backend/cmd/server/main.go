package main

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/organization"
)

func main() {
	http.HandleFunc("/organization/member/", organization.MemberHandler)
	http.HandleFunc("/organization/", organization.Handler)
	http.Handle("/", http.FileServer(http.Dir("./")))

	log.Print("Listening at http://localhost:8443/")
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))
}
