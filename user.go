package interview

import (
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
