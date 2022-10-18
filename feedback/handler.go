package feedback

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__Host-UserUUID")
	if err != nil {
		log.Printf("Error parsing cookie for /feedback: %s", err)
		http.ServeFile(w, r, "./error/unauthenticated.html")
		return
	}

	if cookie.Value == "" {
		log.Printf("Error parsing cookie for /feedback: %s", err)
		http.ServeFile(w, r, "./error/unauthenticated.html")
		return
	}

	http.ServeFile(w, r, "./feedback/index.html")
}
