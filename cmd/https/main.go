package main

import (
	"log"
	"net/http"

	"github.com/justindfuller/interviews/api"
)

func main() {
	api.Handlers()

	/* s := http.Server{
		Addr: ":8443",
		Handler: nil,
		TLSConfig: &tls.Config{

		},
		ReadTimeout: 1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout: 1 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	} */

	log.Print("Listening at https://localhost:8443/")
	log.Fatal(http.ListenAndServeTLS(":8443", "./cmd/https/server.crt", "./cmd/https/server.key", nil))
}
