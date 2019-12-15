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

type PublicContractVersion struct {
	FirstName         string `json:"firstName"`
	SecondName        string `json:"secondName"`
	JobTitle          string `json:"jobTitle"`
	CompanyName       string `json:"companyName"`
	FreelancerGrade   int    `json:"freelancerGrade"`
	FreelancerComment string `json:"freelancerComment"`
	ClientGrade       int    `json:"clientGrade"`
	ClientComment     string `json:"clientComment"`
	Status            string `json:"status"`
}

func (freel *ExtendFreelancer) Sanitize(sanitizer *bluemonday.Policy) {
	freel.F.Sanitize(sanitizer)
	freel.FirstName = sanitizer.Sanitize(freel.FirstName)
	freel.SecondName = sanitizer.Sanitize(freel.SecondName)
}

func (pbc *PublicContractVersion) Sanitize(sanitizer *bluemonday.Policy) {
	pbc.FirstName = sanitizer.Sanitize(pbc.FirstName)
	pbc.SecondName = sanitizer.Sanitize(pbc.SecondName)
	pbc.JobTitle = sanitizer.Sanitize(pbc.JobTitle)
	pbc.CompanyName = sanitizer.Sanitize(pbc.CompanyName)
	pbc.FreelancerComment = sanitizer.Sanitize(pbc.FreelancerComment)
	pbc.ClientComment = sanitizer.Sanitize(pbc.ClientComment)
}
