package company

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Create() (*model.Company, error)
	Find(id int64) (*model.CompanyOutput, error)
	Edit(userId int64, company *model.Company) (*model.CompanyOutput, error)
}