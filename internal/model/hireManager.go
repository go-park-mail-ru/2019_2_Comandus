package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

type HireManager struct {
	ID					int64		`json:"id" valid:"int, optional"`
	AccountID 			int64		`json:"accountId" valid:"int, optional"`
	RegistrationDate	time.Time	`json:"registrationDate" valid:"-l"`
	Location			string 		`json:"location" valid:"-"`
	CompanyID			int64 		`json:"companyId" valid:"int, optional"`
}


func (hireMan *HireManager) Sanitize (sanitizer *bluemonday.Policy)  {
	hireMan.Location = sanitizer.Sanitize(hireMan.Location)
}