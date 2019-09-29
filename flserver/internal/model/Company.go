package model

type Company struct {
	ID			int    `json:"id"`
	CompanyName string `json:"companyName"`
	Site 		string `json:"site"`
	TagLine 	string `json:"tagline"`
	Description string `json:"description"`
	Country 	string `json:"country"`
	City 		string `json:"city"`
	Address 	string `json:"address"`
	Phone 		string `json:"phone"`
}

