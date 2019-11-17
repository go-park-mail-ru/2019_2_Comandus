package company

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Create(c *model.Company) error
	Find(id int64) (*model.Company, error)
	Edit(u *model.User, c *model.Company) error
}