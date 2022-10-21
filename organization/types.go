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

func NewQuestion(text string) (Question, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Question{}, errors.Wrap(err, "error creating ID for Question")
	}

	return Question{
		ID:   id,
		Text: text,
	}, nil
}

type Question struct {
	ID   uuid.UUID
	Text string
}

func NewFeedback(role string, questions []Question) (Feedback, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Feedback{}, errors.Wrap(err, "error creating ID for Feedback")
	}

	return Feedback{
		ID:        id,
		Role:      role,
		Questions: questions,
	}, nil
}

type Feedback struct {
	ID        uuid.UUID
	Role      string
	Questions []Question
}

type Organization struct {
	ID       uuid.UUID
	Domain   string
	Users    []User
	Feedback []Feedback
}

func (o Organization) FeedbackByID(id uuid.UUID) (Feedback, error) {
	for _, feedback := range o.Feedback {
		if feedback.ID == id {
			return feedback, nil
		}
	}

	return Feedback{}, errors.New("feedback not found")
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
