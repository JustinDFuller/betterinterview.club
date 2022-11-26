package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	interview "github.com/justindfuller/interviews"
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

	file, err := os.OpenFile("./organizations.json", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}

	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading organizations.json: %s", err)
		return
	}
	file.Close()

	if err := json.Unmarshal(b, &interview.DefaultOrganizations); err != nil {
		log.Printf("Error decoding organizations.json: %s", err)
		return
	}
	log.Print("Initialized DefaultOrganizations from organizations.json")

	go func() {
		write := func() {
			file, err := os.OpenFile("./organizations.json", os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				log.Printf("Error opening file: %s", err)
				return
			}
			defer file.Close()

			interview.DefaultOrganizations.Mutex.Lock()
			defer interview.DefaultOrganizations.Mutex.Unlock()

			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			if err := enc.Encode(&interview.DefaultOrganizations); err != nil {
				log.Printf("Error encoding to gob: %s", err)
				return
			}
		}

		write()
		for _ = range time.Tick(5 * time.Second) {
			write()
		}
	}()

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
	log.Fatal(http.ListenAndServeTLS(":8443", "./cmd/localhost/server.crt", "./cmd/localhost/server.key", nil))
}
