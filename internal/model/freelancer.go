package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

type Freelancer struct {
	ID                int64       `json:"id"`
	AccountId         int64       `json:"accountId"`
	RegistrationDate  time.Time `json:"registrationDate"`
	Country           string    `json:"country"`
	City              string    `json:"city"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	TagLine           string    `json:"tagline"`
	Overview          string    `json:"overview"`
	ExperienceLevelId int64       `json:"experienceLevelId"`
	SpecialityId      int64       `json:"specialityId,string"`
}


func (freel *Freelancer) Sanitize (sanitizer *bluemonday.Policy)  {
	freel.Country = sanitizer.Sanitize(freel.Country)
	freel.City = sanitizer.Sanitize(freel.City)
	freel.Address = sanitizer.Sanitize(freel.Address)
	freel.Phone = sanitizer.Sanitize(freel.Phone)
	freel.TagLine = sanitizer.Sanitize(freel.TagLine)
	freel.Overview = sanitizer.Sanitize(freel.Overview)
}