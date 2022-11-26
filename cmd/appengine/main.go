package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
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

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating storage client: %s", err)
	}

	object := client.Bucket("betterinterviews").Object("organizations.json")

	rc, err := object.NewReader(ctx)
	if err != nil {
		log.Fatalf("Error creating reader for storage bucket: %s", err)
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		log.Fatalf("Error reading data from storage bucket: %s", err)
	}

	if err := json.Unmarshal(body, &interview.DefaultOrganizations); err != nil {
		log.Fatalf("Error decoding data to DefaultOrganizations: %s", err)
	}

	go func() {
		write := func() {
			file := object.NewWriter(ctx)
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
		for _ = range time.Tick(30 * time.Second) {
			write()
		}
	}()

	api.Handlers()

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Listening at %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
