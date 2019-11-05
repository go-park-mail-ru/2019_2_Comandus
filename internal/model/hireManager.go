package model

import (
	"github.com/microcosm-cc/bluemonday"
	"time"
)

type HireManager struct {
	ID					int64		`json:"id"`
	AccountID 			int64		`json:"accountId"`
	RegistrationDate	time.Time	`json:"registrationDate"`
	Location			string 		`json:"location"`
	CompanyID			int64 		`json:"companyId"`
}


func (hireMan *HireManager) Sanitize (sanitizer *bluemonday.Policy)  {
	sanitizer.Sanitize(hireMan.Location)
}