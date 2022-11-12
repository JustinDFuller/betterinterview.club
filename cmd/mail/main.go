package main

import (
	"log"

	interview "github.com/justindfuller/interviews"
)

func main() {
	to := "justindanielfuller@gmail.com"
	subject := "Log in to Better Interviews"
	html := "<h1>Better Interviews</h1><a href=\"https://localhost:8443\">Log In</a>"

	if err := interview.Email(to, subject, html); err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent subject '%s' to '%s'", subject, to)
}
