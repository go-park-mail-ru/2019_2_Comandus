package model

import "time"

type HireManager struct {
	ID					int		`json:"id"`
	AccountID 			int64		`json:"accountId"`
	RegistrationDate	time.Time	`json:"registrationDate"`
	Location			string 		`json:"location"`
	CompanyID			int 		`json:"companyId"`
}