package organization

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID
	Email string
}

type Organization struct {
	ID     uuid.UUID
	Domain string
	Users  []User
}

type Organizations struct {
	byDomain map[string]Organization
	mutex    sync.Mutex
}

func (orgs *Organizations) Get(domain string) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	org, found := orgs.byDomain[domain]
	if !found {
		return Organization{}, errors.New("organization not found")
	}

	return org, nil
}

func (orgs *Organizations) Add(org Organization) error {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	if _, found := orgs.byDomain[org.Domain]; found {
		return errors.New("organization already exists")
	}

	orgs.byDomain[org.Domain] = org

	return nil
}

var organizations Organizations

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading /organization body: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}
	defer r.Body.Close()

	query, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("Error parsing query from /organization body: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	email, err := mail.ParseAddress(query.Get("email"))
	if err != nil {
		log.Printf("Error parsing email from /organization email parameter: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	parts := strings.Split(email.Address, "@")
	if len(parts) < 2 {
		log.Printf("Error parsing email from /organization email parameter: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	orgID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error creating org UUID in /organization: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error creating user UUID in /organization: %s", err)
		http.ServeFile(w, r, "./error/index.html")
		return
	}

	org := Organization{
		ID:     orgID,
		Domain: parts[1],
		Users: []User{
			{
				ID:    userID,
				Email: email.Address,
			},
		},
	}
	if err := organizations.Add(org); err != nil {
		log.Printf("Error adding organization '%s' from /organization: %s", org.Domain, err)
		http.ServeFile(w, r, "./organization/exists.html")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-UserUUID",
		Value:    userID.String(),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 31), // One month
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	log.Printf("New Organization: %s %#v", org.Domain, org)
	http.ServeFile(w, r, "./organization/index.html")
}
