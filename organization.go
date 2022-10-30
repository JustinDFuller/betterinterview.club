package interview

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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
