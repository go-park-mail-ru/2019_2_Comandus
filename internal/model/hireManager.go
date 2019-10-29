package model

import "time"

type HireManager struct {
	ID					int		`json:"id" valid:"int, optional"`
	AccountID 			int64		`json:"accountId" valid:"int, optional"`
	RegistrationDate	time.Time	`json:"registrationDate" valid:"-l"`
	Location			string 		`json:"location" valid:"-"`
	CompanyID			int 		`json:"companyId" valid:"int, optional"`
}