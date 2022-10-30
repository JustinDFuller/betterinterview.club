package interview

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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
