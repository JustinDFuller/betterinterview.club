package interview

import (
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrOrgNotFound = errors.New("organization not found")
var ErrUserNotFound = errors.New("user not found")

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
		return Organization{}, ErrOrgNotFound
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

func (orgs *Organizations) FindByUserID(id string) (Organization, User, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	for _, organization := range orgs.byDomain {
		if user, err := organization.FindUserByID(id); err == nil {
			return organization, user, nil
		}
	}

	return Organization{}, User{}, errors.Errorf("organization not found for user ID: %s", id)
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

	return Organization{}, ErrOrgNotFound
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

func (orgs *Organizations) SetFeedback(org Organization, f Feedback) error {
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
			org.Feedback[i] = f
		}
	}

	orgs.byDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddEmailLoginCallback(org Organization, u User) (string, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	// For this one, we must assume it exists
	if orgs.byDomain == nil {
		return "", ErrOrgNotFound
	}

	org, ok := orgs.byDomain[org.Domain]
	if !ok {
		return "", ErrOrgNotFound
	}

	cbID, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "error creating callback ID")
	}

	var found bool
	for index := range org.Users {
		user := org.Users[index]

		if u.ID == user.ID {
			org.Users[index].CallbackID = cbID
			found = true
		}
	}

	if !found {
		return "", errors.New("user not found")
	}

	orgs.byDomain[org.Domain] = org

	return cbID.String(), nil
}

func (orgs *Organizations) FindOrCreateByEmail(email string) (Organization, User, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		orgs.byDomain = map[string]Organization{}
	}

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return Organization{}, User{}, errors.Errorf("invalid email address: %s", email)
	}

	org, ok := orgs.byDomain[parts[1]]
	if !ok {
		o, err := NewOrganization(parts[1])
		if err != nil {
			return Organization{}, User{}, errors.Wrap(err, "error creating new organization")
		}

		org = o
	}

	user, err := org.FindUserByEmail(email)
	if err != nil {
		user, err = NewUser(email)
		if err != nil {
			return Organization{}, User{}, errors.Wrap(err, "error creating new user")
		}
		org.Users = append(org.Users, user)
	}

	orgs.byDomain[parts[1]] = org

	return org, user, nil
}

func (orgs *Organizations) FindEmailLoginCallback(id string) (User, error) {
	orgs.mutex.Lock()
	defer orgs.mutex.Unlock()

	if orgs.byDomain == nil {
		return User{}, ErrOrgNotFound
	}

	for _, org := range orgs.byDomain {
		for userI := range org.Users {
			user := org.Users[userI]

			if user.CallbackID.String() == id {
				org.Users[userI].CallbackID = uuid.UUID{}
				orgs.byDomain[org.Domain] = org

				return user, nil
			}
		}
	}

	return User{}, ErrUserNotFound
}
