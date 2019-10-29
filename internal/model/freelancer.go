package model

import "time"

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
