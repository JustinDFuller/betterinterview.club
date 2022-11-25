package interview

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FeedbackRequest struct {
	ID                uuid.UUID
	CreatedAt         time.Time
	CandidateName     string
	InterviewerEmails []string
	Responses         []FeedbackResponse
}

func NewFeedbackRequest(candidate string, emails ...string) (FeedbackRequest, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return FeedbackRequest{}, errors.Wrap(err, "error creating ID for FeedbackRequest")
	}

	var interviewerEmails []string
	emailMap := map[string]bool{}

	for _, email := range emails {
		// No extra space, case insensitive
		email := strings.TrimSpace(strings.ToLower(email))

		// No duplicate emails
		if emailMap[email] {
			continue
		}

		// No empty emails
		if email == "" {
			continue
		}

		interviewerEmails = append(interviewerEmails, email)
		emailMap[email] = true
	}

	return FeedbackRequest{
		ID:                id,
		CreatedAt:         time.Now(),
		CandidateName:     candidate,
		InterviewerEmails: interviewerEmails,
	}, nil
}
