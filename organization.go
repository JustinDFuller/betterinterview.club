package interview

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Test
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

func (o Organization) FeedbackByRequestID(id uuid.UUID) (Feedback, FeedbackRequest, error) {
	for _, feedback := range o.Feedback {
		for _, request := range feedback.Requests {
			if request.ID == id {
				return feedback, request, nil
			}
		}
	}

	return Feedback{}, FeedbackRequest{}, errors.New("feedback request not found")
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

func (org *Organization) IsDifferentDomain(email string) (bool, error) {
	if len(email) == 0 {
		return false, errors.New("invalid email")
	}

	split := strings.Split(email, "@")

	if len(split) != 2 {
		return false, errors.New("invalid email")
	}

	domain := strings.TrimSpace(strings.ToLower(split[1]))
	if len(domain) == 0 {
		return false, errors.New("invalid email")
	}

	return domain != strings.TrimSpace(strings.ToLower(org.Domain)), nil
}

func NewOrganization(domain string) (Organization, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Organization{}, errors.Wrap(err, "error creating ID for Organization")
	}

	return Organization{
		ID:     id,
		Domain: domain,
	}, nil
}
