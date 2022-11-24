package interview

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FeedbackRequest struct {
	ID                uuid.UUID
	CandidateName     string
	InterviewerEmails []string
}

func NewFeedbackRequest(candidate string, emails ...string) (FeedbackRequest, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return FeedbackRequest{}, errors.Wrap(err, "error creating ID for FeedbackRequest")
	}

	var interviewerEmails []string
	for _, email := range emails {
		if email != "" {
			interviewerEmails = append(interviewerEmails, email)
		}
	}

	return FeedbackRequest{
		ID:                id,
		CandidateName:     candidate,
		InterviewerEmails: interviewerEmails,
	}, nil
}
