package interview

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var OrgNotFound = errors.New("organization not found")

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
		return Organization{}, OrgNotFound
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

	o, found := orgs.byDomain[org.Domain]
	if !found {
		return o, errors.New("organization does not exist")
	}

	for _, user := range o.Users {
		if user.Email == u.Email {
			return o, errors.New("user already exists")
		}
	}

	o.Users = append(o.Users, u)
	orgs.byDomain[o.Domain] = o

	return o, nil
}

func (orgs *Organizations) FindByUserEmail(email string) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	for _, organization := range orgs.byDomain {
		if _, err := organization.FindUserByEmail(email); err == nil {
			return organization, nil
		}
	}

	return Organization{}, errors.Errorf("organization not found for user email: %s", email)
}

func (orgs *Organizations) FindByUserID(id string) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	for _, organization := range orgs.byDomain {
		if _, err := organization.FindUserByID(id); err == nil {
			return organization, nil
		}
	}

	return Organization{}, errors.Errorf("organization not found for user ID: %s", id)
}

func (orgs *Organizations) FindByDomain(email string) (Organization, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return Organization{}, errors.Errorf("invalid email address: %s", email)
	}

	for _, organization := range orgs.byDomain {
		if organization.Domain == parts[1] {
			return organization, nil
		}
	}

	return Organization{}, errors.Errorf("organization not found for domain: %s", parts[1])
}

func (orgs *Organizations) AddFeedback(org Organization, f Feedback) error {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	org, found := orgs.byDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	org.Feedback = append(org.Feedback, f)
	orgs.byDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddFeedbackResponse(org Organization, f Feedback, fr FeedbackResponse) error {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	org, found := orgs.byDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	for i, feedback := range org.Feedback {
		if feedback.ID == f.ID {
			org.Feedback[i].Responses = append(feedback.Responses, fr)
		}
	}

	orgs.byDomain[org.Domain] = org

	return nil
}
