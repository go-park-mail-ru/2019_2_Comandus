package company

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Find(id int64) (*model.Company, error)
	Edit(u *model.User, company *model.Company) error
}