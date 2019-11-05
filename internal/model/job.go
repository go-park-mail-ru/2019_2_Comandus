package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

const (
	JobStateCreate = "created"
	JobStateFound = "found"
	JobStateClosed = "closed"
	)

type Job struct {
	ID                int64     `json:"id"`
	HireManagerId     int64     `json:"hireManagerId"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Files             string  `json:"files"`
	SpecialityId      int64     `json:"specialityId,string"`
	ExperienceLevelId int64     `json:"experienceLevelId,string"`
	PaymentAmount     float64 `json:"paymentAmount,string"`
	Country           string  `json:"country"`
	City              string  `json:"city"`
	JobTypeId         int64     `json:"jobTypeId,string"`
	Date			  time.Time `json:"date"`
	Status			  string `json:"status,string"`
}

func (j *Job) IsEqual(job Job) bool {
	return j.ID == job.ID &&
		j.HireManagerId == job.HireManagerId &&
		j.Title == job.Title &&
		j.Description == job.Description &&
		j.Files == job.Files &&
		j.ExperienceLevelId == job.ExperienceLevelId &&
		j.PaymentAmount == job.PaymentAmount &&
		j.Country == job.Country &&
		j.City == job.City &&
		j.JobTypeId == job.JobTypeId &&
		j.Date == job.Date &&
		j.Status == job.Status
}

func (j *Job) BeforeCreate() {
	j.Date = time.Now()
	j.Status = JobStateCreate
}

func (job *Job) Sanitize (sanitizer *bluemonday.Policy)  {
	sanitizer.Sanitize(job.Title)
	sanitizer.Sanitize(job.Description)
	sanitizer.Sanitize(job.Files)
	sanitizer.Sanitize(job.Country)
	sanitizer.Sanitize(job.City)
}