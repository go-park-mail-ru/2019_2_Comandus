package model

type Job struct {
	ID                int     `json:"id"`
	HireManagerId     int     `json:"hireManagerId,string"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Files             string  `json:"files"`
	SpecialityId      int     `json:"specialityId,string"`
	ExperienceLevelId int     `json:"experienceLevelId,string"`
	PaymentAmout      float64 `json:"paymentAmount,string"`
	Country           string  `json:"country"`
	City              string  `json:"city"`
	JobTypeId         int     `json:"jobTypeId,string"`
}
