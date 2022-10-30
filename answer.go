package interview

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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
