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
	sanitizer.Sanitize(comp.CompanyName)
	sanitizer.Sanitize(comp.Site)
	sanitizer.Sanitize(comp.TagLine)
	sanitizer.Sanitize(comp.Description)
	sanitizer.Sanitize(comp.Country)
	sanitizer.Sanitize(comp.City)
	sanitizer.Sanitize(comp.Address)
	sanitizer.Sanitize(comp.Phone)
}