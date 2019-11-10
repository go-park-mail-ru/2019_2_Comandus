package company

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(company *model.Company) error
	Find(int64) (*model.Company, error)
	Edit(company *model.Company) error
}