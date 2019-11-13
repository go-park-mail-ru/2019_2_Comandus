package model

import (
	"github.com/microcosm-cc/bluemonday"
)

type HireManager struct {
	ID					int64		`json:"id" valid:"int, optional"`
	AccountID 			int64		`json:"accountId" valid:"int, optional"`
	Location			string 		`json:"location" valid:"-"`
	CompanyID			int64 		`json:"companyId" valid:"int, optional"`
}


func (hireMan *HireManager) Sanitize (sanitizer *bluemonday.Policy)  {
	hireMan.Location = sanitizer.Sanitize(hireMan.Location)
}