package auth

import (
	"html/template"
	"log"
	"net/http"
)

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("email.html").ParseFiles("./auth/email.html", "index.css")
	if err != nil {
		log.Printf("Error parsing template for /auth/email/: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	if err := t.Execute(w, nil); err != nil {
		log.Printf("Error executing template for /auth/email: %s", err)
	}
	return
}
