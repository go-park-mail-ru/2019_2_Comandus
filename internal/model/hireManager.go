package model

import "time"

type HireManager struct {
	ID					int64		`json:"id"`
	AccountID 			int64		`json:"accountId"`
	RegistrationDate	time.Time	`json:"registrationDate"`
	Location			string 		`json:"location"`
	CompanyID			int64 		`json:"companyId"`
}