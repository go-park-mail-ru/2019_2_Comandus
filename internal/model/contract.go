package model

import (
	"errors"
	"time"
)

const (
	ContractStatusDone = "done"
	ContractStatusUnderDevelopment = "develop"
	ContractStatusCanceled = "cancel"
	ContractStatusReviewed = "reviewed"
	ContractMinGrade = 0
	ContractMaxGrade = 5
)

type Contract struct {
	ID				int64		`json:"id"`
	ResponseID		int64 		`json:"responseId"`
	CompanyID		int64		`json:"companyId"`
	FreelancerID	int64		`json:"freelancerId"`
	StartTime		time.Time	`json:"startTime"`
	EndTime			time.Time	`json:"endTime"`
	Status			string		`json:"status"`
	Grade			int			`json:"grade"`
	PaymentAmount	float64		`json:"paymentAmount"`
}


func (c * Contract) Validate(lastId int64) error {
	if c.ID != lastId {
		return errors.New("cant change contract ID")
	}

	if c.ResponseID == 0 || c.CompanyID == 0 || c.FreelancerID == 0 {
		return errors.New("responseId, companyId and freelancerId cant be null")
	}

	if len(c.Status) == 0 {
		return errors.New("wrong status")
	}
	return nil
}