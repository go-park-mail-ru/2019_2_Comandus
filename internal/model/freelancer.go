package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

type Freelancer struct {
	ID                int64       `json:"id" valid:"int, optional"`
	AccountId         int64       `json:"accountId" valid:"int, optional"`
	RegistrationDate  time.Time `json:"registrationDate" valid:"-"`
	Country           string    `json:"country" valid:"utfletter"`
	City              string    `json:"city" valid:"utfletter"`
	Address           string    `json:"address" valid:"-"`
	Phone             string    `json:"phone" valid:"-"`
	TagLine           string    `json:"tagline" valid:"-"`
	Overview          string    `json:"overview" valid:"-"`
	ExperienceLevelId int       `json:"experienceLevelId" valid:"in(1|2|3)"`
	SpecialityId      int       `json:"specialityId,string" valid:"int"`
}


func (freel *Freelancer) Sanitize (sanitizer *bluemonday.Policy)  {
	freel.Country = sanitizer.Sanitize(freel.Country)
	freel.City = sanitizer.Sanitize(freel.City)
	freel.Address = sanitizer.Sanitize(freel.Address)
	freel.Phone = sanitizer.Sanitize(freel.Phone)
	freel.TagLine = sanitizer.Sanitize(freel.TagLine)
	freel.Overview = sanitizer.Sanitize(freel.Overview)
}