package interview

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewFeedbackResponse(creatorID uuid.UUID, answers []Answer, recommend bool) (FeedbackResponse, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return FeedbackResponse{}, errors.Wrap(err, "error creating ID for FeedbackResponse")
	}

	return FeedbackResponse{
		ID:        id,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Answers:   answers,
		Recommend: recommend,
	}, nil
}

type FeedbackResponse struct {
	ID        uuid.UUID
	CreatorID uuid.UUID
	CreatedAt time.Time
	Answers   []Answer
	Recommend bool
}
