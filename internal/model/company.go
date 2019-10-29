package model

type Company struct {
	ID			int    `json:"id" valid:"int , optional"`
	CompanyName string `json:"companyName" valid:"utfletternum, required"`
	Site 		string `json:"site" valid:"url"`
	TagLine 	string `json:"tagline" valid:"- , optional"`
	Description string `json:"description" valid:"-"`
	Country 	string `json:"country" valid:"utfletter"`
	City 		string `json:"city" valid:"utfletter"`
	Address 	string `json:"address" valid:"-"`
	Phone 		string `json:"phone" valid:"- , optional"`
}

