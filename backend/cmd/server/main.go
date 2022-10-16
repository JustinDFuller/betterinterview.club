package main

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/organization"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/organization", organization.Handler)

	log.Print("Listening at http://localhost:8443/")
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))
}
