package interview

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewUser(email string) (User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return User{}, errors.Wrap(err, "error creating ID for User")
	}

	return User{
		ID:    id,
		Email: email,
	}, nil
}

type User struct {
	ID    uuid.UUID
	Email string
}

func NewAnswer(questionID string, response bool) (Answer, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Answer{}, errors.Wrap(err, "error creating ID for Answer")
	}

	qID, err := uuid.Parse(questionID)
	if err != nil {
		return Answer{}, errors.Wrap(err, "error creating QuestionID for Answer")
	}

	return Answer{
		ID:         id,
		QuestionID: qID,
		Response:   response,
	}, nil
}

type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Response   bool
}

func NewFeedbackResponse(creatorID uuid.UUID, answers []Answer) (FeedbackResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return FeedbackResponse{}, errors.Wrap(err, "error creating ID for FeedbackResponse")
	}

	return FeedbackResponse{
		ID:        id,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Answers:   answers,
	}, nil
}

type FeedbackResponse struct {
	ID        uuid.UUID
	CreatorID uuid.UUID
	CreatedAt time.Time
	Answers   []Answer
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

func NewFeedback(creatorID uuid.UUID, role string, questions []Question) (Feedback, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Feedback{}, errors.Wrap(err, "error creating ID for Feedback")
	}

	return Feedback{
		ID:        id,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Role:      role,
		Questions: questions,
	}, nil
}

type Feedback struct {
	ID        uuid.UUID
	CreatorID uuid.UUID
	CreatedAt time.Time
	Role      string
	Questions []Question
	Responses []FeedbackResponse
}

func (f *Feedback) QuestionByID(id string) (Question, error) {
	for _, question := range f.Questions {
		if question.ID.String() == id {
			return question, nil
		}
	}

	return Question{}, errors.Errorf("question not found: %s", id)
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

func (org *Organization) FindUserByEmail(email string) (User, error) {
	for _, user := range org.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.Errorf("User not found for user email: %s", email)
}

func (org *Organization) FindUserByID(id string) (User, error) {
	for _, user := range org.Users {
		if user.ID.String() == id {
			return user, nil
		}
	}

	return User{}, errors.Errorf("User not found for user ID: %s", id)
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
