package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/justindfuller/interviews/api"
)

func main() {
	f, err := os.ReadFile("./.secrets.env")
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range strings.Split(string(f), "\n") {
		// Last line is just a newline
		if line == "" {
			continue
		}

		parts := strings.Split(line, "export ")
		parts = strings.Split(parts[1], "=")
		if err := os.Setenv(parts[0], parts[1]); err != nil {
			log.Fatal(err)
		}
	}

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
