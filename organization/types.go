package organization

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (orgs *Organizations) AddUser(org Organization, u User) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	o := orgs.byDomain[org.Domain]
	o.Users = append(o.Users, u)
	orgs.byDomain[o.Domain] = o

	return o, nil
}

func (orgs *Organizations) FindByUserID(id string) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	var org Organization
	for _, organization := range orgs.byDomain {
		for _, user := range organization.Users {
			if user.ID.String() == id {
				return organization, nil
			}
		}
	}

	return org, errors.Errorf("organization not found for user ID: %s", id)
}
