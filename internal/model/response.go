package model

import (
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"time"
)

const (
	ResponseStatusBlock = "block"
	ResponseStatusReview = "review"
	ResponseStatusDenied = "denied"
	ResponseStatusAccepted = "accepted"
	)

type Response struct {
	ID					int64     `json:"id"`
	FreelancerId		int64     `json:"freelancerId"`
	JobId				int64     `json:"jobId"`
	Files				string    `json:"files,string"`
	Date				time.Time `json:"date"`
	StatusManager		string    `json:"statusManager,string"`
	StatusFreelancer	string    `json:"statusFreelancer,string"`
	PaymentAmount 		float64	  `json:"PaymentAmount,string"`
}

func (r *Response) BeforeCreate() {
	r.StatusManager = ResponseStatusReview
	r.StatusFreelancer = ResponseStatusBlock
	r.Date = time.Now()
}

// validation before create and edit
// for create lastID = 0
func (r * Response) Validate(lastID int64) error {
	if !(r.StatusManager == ResponseStatusReview ||
		r.StatusManager == ResponseStatusDenied ||
		r.StatusManager == ResponseStatusAccepted) {
		return errors.New("wrong manager response status")
	}
	if !(r.StatusFreelancer == ResponseStatusReview ||
		r.StatusFreelancer == ResponseStatusDenied ||
		r.StatusFreelancer == ResponseStatusAccepted ||
		r.StatusFreelancer == ResponseStatusBlock) {
		return errors.New("wrong freelancer response status")
	}
	if r.Date.IsZero() {
		return errors.New("wrong date")
	}
	if r.ID != lastID {
		return errors.New("current id does not match last id")
	}
	if r.FreelancerId == 0 || r.JobId == 0 {
		return errors.New("wrong relationships between tables")
	}
	return nil
}

func (r * Response) IsEqual(response *Response) bool {
	return r.ID == response.ID &&
		r.FreelancerId == response.FreelancerId &&
		r.JobId == response.JobId &&
		r.Files == response.Files &&
		r.Date == response.Date &&
		r.StatusManager == response.StatusManager &&
		r.StatusFreelancer == response.StatusFreelancer
}


func (resp *Response) Sanitize (sanitizer *bluemonday.Policy)  {
	resp.Files = sanitizer.Sanitize(resp.Files)
}