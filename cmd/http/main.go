package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justindfuller/interviews/api"
)

func main() {
	api.Handlers()

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Listening at %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
