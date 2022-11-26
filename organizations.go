package interview

import (
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var DefaultOrganizations = &Organizations{}

var ErrOrgNotFound = errors.New("organization not found")
var ErrUserNotFound = errors.New("user not found")

type Organizations struct {
	ByDomain map[string]Organization
	Mutex    sync.Mutex `json:"-"`
}

func (orgs *Organizations) Get(domain string) (Organization, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	org, found := orgs.ByDomain[domain]
	if !found {
		return Organization{}, ErrOrgNotFound
	}

	return org, nil
}

func (orgs *Organizations) Add(org Organization) error {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	if _, found := orgs.ByDomain[org.Domain]; found {
		return errors.New("organization already exists")
	}

	orgs.ByDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddUser(org Organization, u User) (Organization, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	o, found := orgs.ByDomain[org.Domain]
	if !found {
		return o, errors.New("organization does not exist")
	}

	for _, user := range o.Users {
		if user.Email == u.Email {
			return o, errors.New("user already exists")
		}
	}

	o.Users = append(o.Users, u)
	orgs.ByDomain[o.Domain] = o

	return o, nil
}

func (orgs *Organizations) FindByUserEmail(email string) (Organization, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	for _, organization := range orgs.ByDomain {
		if _, err := organization.FindUserByEmail(email); err == nil {
			return organization, nil
		}
	}

	return Organization{}, errors.Errorf("organization not found for user email: %s", email)
}

func (orgs *Organizations) FindByUserID(id string) (Organization, User, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	for _, organization := range orgs.ByDomain {
		if user, err := organization.FindUserByID(id); err == nil {
			return organization, user, nil
		}
	}

	return Organization{}, User{}, errors.Errorf("organization not found for user ID: %s", id)
}

func (orgs *Organizations) FindByDomain(email string) (Organization, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return Organization{}, errors.Errorf("invalid email address: %s", email)
	}

	for _, organization := range orgs.ByDomain {
		if organization.Domain == parts[1] {
			return organization, nil
		}
	}

	return Organization{}, ErrOrgNotFound
}

func (orgs *Organizations) AddFeedback(org Organization, f Feedback) error {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	org, found := orgs.ByDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	org.Feedback = append(org.Feedback, f)
	orgs.ByDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddFeedbackRequest(org Organization, f Feedback, request FeedbackRequest) error {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	org, found := orgs.ByDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	for i, feedback := range org.Feedback {
		if feedback.ID == f.ID {
			org.Feedback[i].Requests = append(feedback.Requests, request)
		}
	}

	orgs.ByDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddFeedbackResponse(org Organization, request FeedbackRequest, response FeedbackResponse) error {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	org, found := orgs.ByDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	for i, feedback := range org.Feedback {
		for j, feedbackRequest := range feedback.Requests {
			if feedbackRequest.ID == request.ID {
				org.Feedback[i].Requests[j].Responses = append(org.Feedback[i].Requests[j].Responses, response)
				orgs.ByDomain[org.Domain] = org
				return nil
			}
		}
	}

	return errors.New("feedback request does not exist")
}

func (orgs *Organizations) SetFeedback(org Organization, f Feedback) error {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	org, found := orgs.ByDomain[org.Domain]
	if !found {
		return errors.New("organization does not exist")
	}

	for i, feedback := range org.Feedback {
		if feedback.ID == f.ID {
			org.Feedback[i] = f
		}
	}

	orgs.ByDomain[org.Domain] = org

	return nil
}

func (orgs *Organizations) AddEmailLoginCallback(org Organization, u User) (string, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	// For this one, we must assume it exists
	if orgs.ByDomain == nil {
		return "", ErrOrgNotFound
	}

	org, ok := orgs.ByDomain[org.Domain]
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

	orgs.ByDomain[org.Domain] = org

	return cbID.String(), nil
}

func (orgs *Organizations) FindOrCreateByEmail(email string) (Organization, User, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		orgs.ByDomain = map[string]Organization{}
	}

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return Organization{}, User{}, errors.Errorf("invalid email address: %s", email)
	}

	org, ok := orgs.ByDomain[parts[1]]
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

	orgs.ByDomain[org.Domain] = org

	return org, user, nil
}

func (orgs *Organizations) FindEmailLoginCallback(id string) (User, error) {
	orgs.Mutex.Lock()
	defer orgs.Mutex.Unlock()

	if orgs.ByDomain == nil {
		return User{}, ErrOrgNotFound
	}

	for _, org := range orgs.ByDomain {
		for userI := range org.Users {
			user := org.Users[userI]

			if user.CallbackID.String() == id {
				org.Users[userI].CallbackID = uuid.UUID{}
				orgs.ByDomain[org.Domain] = org

				return user, nil
			}
		}
	}

	return User{}, ErrUserNotFound
}
