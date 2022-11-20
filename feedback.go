package interview

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewFeedback(creatorID uuid.UUID, team, role string, emails []string, questions []Question) (Feedback, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Feedback{}, errors.Wrap(err, "error creating ID for Feedback")
	}

	return Feedback{
		ID:        id,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Team:      team,
		Role:      role,
		Emails:    emails,
		Questions: questions,
	}, nil
}

type Feedback struct {
	ID          uuid.UUID
	CreatorID   uuid.UUID
	CreatedAt   time.Time
	Team        string
	Role        string
	Emails      []string
	Questions   []Question
	Responses   []FeedbackResponse
	Closed      bool
	CloseReason string
}

func (f *Feedback) QuestionByID(id string) (Question, error) {
	for _, question := range f.Questions {
		if question.ID.String() == id {
			return question, nil
		}
	}

	return Question{}, errors.Errorf("question not found: %s", id)
}
