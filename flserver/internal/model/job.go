package model

type Job struct {
	ID int `json:"id"`
	HireManagerId int `json:"hireManagerId"`
	Title string `json:"title"`
	Description string `json:"description"`
	Files string `json:"files"`
	SpecialityId int `json:"specialityId"`
	ExperienceLevelId int `json:"experienceLevelId"`
	PaymentAmout float64 `json:"paymentAmount"`
	Country string `json:"country"`
	City string `json:"city"`
	JobTypeId int `json:"jobTypeId"`
}
