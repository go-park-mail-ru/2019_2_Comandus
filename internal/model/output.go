package model

import "github.com/microcosm-cc/bluemonday"

type ExtendFreelancer struct {
	F          *Freelancer `json:"freelancer"`
	FirstName  string      `json:"firstName"`
	SecondName string      `json:"secondName"`
}

type Review struct {
	CompanyName   string `json:"companyName"`
	JobTitle      string `json:"jobTitle"`
	ClientGrade   int    `json:"clientGrade"`
	ClientComment string `json:"clientComment"`
}

// what for
type OutputResponse struct {
	Id int64 `json:"id"`
}

//
type ExtendResponse struct {
	R          *Response `json:"Response"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"lastName"`
	JobTitle   string    `json:"jobTitle, string"`
}


type ExtendedOutputFreelancer struct {
	OuFreel    *FreelancerOutput `json:"freelancer"`
	FirstName  string            `json:"firstName"`
	SecondName string            `json:"secondName"`
}


func (freel *ExtendFreelancer) Sanitize(sanitizer *bluemonday.Policy) {
	freel.F.Sanitize(sanitizer)
	freel.FirstName = sanitizer.Sanitize(freel.FirstName)
	freel.SecondName = sanitizer.Sanitize(freel.SecondName)
}