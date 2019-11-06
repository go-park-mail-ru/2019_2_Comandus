package model

import "github.com/microcosm-cc/bluemonday"

type Company struct {
	ID			int64    `json:"id"`
	CompanyName string `json:"companyName"`
	Site 		string `json:"site"`
	TagLine 	string `json:"tagline"`
	Description string `json:"description"`
	Country 	string `json:"country"`
	City 		string `json:"city"`
	Address 	string `json:"address"`
	Phone 		string `json:"phone"`
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