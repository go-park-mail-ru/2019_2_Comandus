package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

const (
	JobStateCreate = "created"
	JobStateFound = "found"
	JobStateClosed = "closed"
	JobStateDeleted = "deleted"
	)

type Job struct {
	ID                int64     `json:"id" valid:"int, optional"`
	HireManagerId     int64   `json:"hireManagerId,string" valid:"int, optional"`
	Title             string  `json:"title" valid:"utfletternum, required"`
	Description       string  `json:"description"valid:"- , optional"`
	Files             string  `json:"files" valid:"-"`
	SpecialityId      int64     `json:"specialityId,string" valid:"int, optional"`
	ExperienceLevelId int64     `json:"experienceLevelId,string" valid:"in(1|2|3)"`
	PaymentAmount     float32 `json:"paymentAmount,string" valid:"float"`
	Country           int64  `json:"country" valid:"utfletternum, optional"`
	City              int64  `json:"city" valid:"utfletternum, optional"`
	JobTypeId         int64     `json:"jobTypeId,string" valid:"int, optional"`
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

func (j *Job) Sanitize (sanitizer *bluemonday.Policy)  {
	j.Title = sanitizer.Sanitize(j.Title)
	j.Description = sanitizer.Sanitize(j.Description)
	j.Files = sanitizer.Sanitize(j.Files)
}