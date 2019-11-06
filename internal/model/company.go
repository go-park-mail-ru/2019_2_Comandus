package model

import "github.com/microcosm-cc/bluemonday"

type Company struct {
	ID			int64    `json:"id" valid:"int , optional"`
	CompanyName string `json:"companyName" valid:"utfletternum, required"`
	Site 		string `json:"site" valid:"url"`
	TagLine 	string `json:"tagline" valid:"- , optional"`
	Description string `json:"description" valid:"-"`
	Country 	string `json:"country" valid:"utfletter"`
	City 		string `json:"city" valid:"utfletter"`
	Address 	string `json:"address" valid:"-"`
	Phone 		string `json:"phone" valid:"- , optional"`
}

func (comp *Company) Sanitize (sanitizer *bluemonday.Policy)  {
	comp.CompanyName = sanitizer.Sanitize(comp.CompanyName)
	comp.Site = sanitizer.Sanitize(comp.Site)
	comp.TagLine = sanitizer.Sanitize(comp.TagLine)
	comp.Description = sanitizer.Sanitize(comp.Description)
	comp.Country = sanitizer.Sanitize(comp.Country)
	comp.City = sanitizer.Sanitize(comp.City)
	comp.Address = sanitizer.Sanitize(comp.Address)
	comp.Phone = sanitizer.Sanitize(comp.Phone)
}