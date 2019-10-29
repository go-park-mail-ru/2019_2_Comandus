package model

import "time"

type Freelancer struct {
	ID                int       `json:"id"`
	AccountId         int64       `json:"accountId"`
	RegistrationDate  time.Time `json:"registrationDate"`
	Country           string    `json:"country"`
	City              string    `json:"city"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	TagLine           string    `json:"tagline"`
	Overview          string    `json:"overview"`
	ExperienceLevelId int       `json:"experienceLevelId"`
	SpecialityId      int       `json:"specialityId,string"`
}
