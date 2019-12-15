package model

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc/job_grpc"
	"time"
)

const (
	ContractStatusDone             = "closed"
	ContractStatusUnderDevelopment = "active"
	ContractStatusDenied           = "denied"
	ContractStatusExpected         = "expected"

	ContractMinGrade = 1
	ContractMaxGrade = 5

	FreelancerReady     = "Ready"
	FreelacncerNotReady = "NotReady"
)

type Contract struct {
	ID                   int64     `json:"id"`
	ResponseID           int64     `json:"responseId"`
	CompanyID            int64     `json:"companyId"`
	FreelancerID         int64     `json:"freelancerId"`
	StartTime            time.Time `json:"startTime"`
	EndTime              time.Time `json:"endTime"`
	Status               string    `json:"status,string"`
	StatusFreelancerWork string    `json:"statusFreelancerWork,string"`
	FreelancerGrade      int       `json:"freelancerGrade"`
	FreelancerComment    string    `json:"freelancerComment,string"`
	ClientGrade          int       `json:"clientGrade"`
	ClientComment        string    `json:"clientComment,string"`
	PaymentAmount        float32   `json:"paymentAmount"`
	TimeEstimation       int       `json:"timeEstimation"`
}

type ContractOutput struct {
	Company    company_grpc.CompanyOutput
	Freelancer freelancer_grpc.ExtendedFreelancer
	Job        job_grpc.Job
	Contract   Contract
}

func (c *Contract) Validate(lastId int64) error {
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
