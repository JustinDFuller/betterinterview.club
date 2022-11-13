package main

import (
	"log"
	"os"
	"strings"

	interview "github.com/justindfuller/interviews"
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

	to := "justindanielfuller@gmail.com"
	subject := "Log in to Better Interviews"
	html := "<h1>Better Interviews</h1><a href=\"https://localhost:8443\">Log In</a>"

	if err := interview.Email(interview.EmailOptions{to, subject, html}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent subject '%s' to '%s'", subject, to)
}
