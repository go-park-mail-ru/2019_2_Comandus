package model

import "github.com/microcosm-cc/bluemonday"

type Company struct {
	ID          int64  `json:"id" valid:"int , optional"`
	CompanyName string `json:"companyName" valid:"utfletternum, required"`
	CompanyOwner string `json:"companyOwner"`
	Site        string `json:"site" valid:"url"`
	TagLine     string `json:"tagline" valid:"- , optional"`
	Description string `json:"description" valid:"-"`
	Country     int64  `json:"country" valid:"utfletter"`
	City        int64  `json:"city" valid:"utfletter"`
	Address     string `json:"address" valid:"-"`
	Phone       string `json:"phone" valid:"- , optional"`
}

type CompanyOutput struct {
	ID          int64
	CompanyName string
	Site        string
	TagLine     string
	Description string
	Country     string
	City        string
	Address     string
	Phone       string
}

func (comp *Company) Sanitize(sanitizer *bluemonday.Policy) {
	comp.CompanyName = sanitizer.Sanitize(comp.CompanyName)
	comp.Site = sanitizer.Sanitize(comp.Site)
	comp.TagLine = sanitizer.Sanitize(comp.TagLine)
	comp.Description = sanitizer.Sanitize(comp.Description)
	comp.Address = sanitizer.Sanitize(comp.Address)
	comp.Phone = sanitizer.Sanitize(comp.Phone)
}

func (comp *CompanyOutput) Sanitize(sanitizer *bluemonday.Policy) {
	comp.CompanyName = sanitizer.Sanitize(comp.CompanyName)
	comp.Site = sanitizer.Sanitize(comp.Site)
	comp.TagLine = sanitizer.Sanitize(comp.TagLine)
	comp.Description = sanitizer.Sanitize(comp.Description)
	comp.Country = sanitizer.Sanitize(comp.Country)
	comp.City = sanitizer.Sanitize(comp.City)
	comp.Address = sanitizer.Sanitize(comp.Address)
	comp.Phone = sanitizer.Sanitize(comp.Phone)
}
