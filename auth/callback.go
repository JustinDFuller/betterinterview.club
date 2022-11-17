package auth

import (
	"log"
	"net/http"
	"time"

	interview "github.com/justindfuller/interviews"
)

func CallbackHandler(organizations *interview.Organizations) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := organizations.FindEmailLoginCallback(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("Error adding email login callback in /auth/login: %s", err)
			http.ServeFile(w, r, "./error/index.html")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "__Host-UserUUID",
			Value:    user.ID.String(),
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 31), // One month
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if redirect := r.URL.Query().Get("redirect"); redirect != "" {
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/organization/", http.StatusSeeOther)
	}
}
