package main

import (
	"fmt"
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

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Listening at %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
