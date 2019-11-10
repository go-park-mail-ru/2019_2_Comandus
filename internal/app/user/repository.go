package user

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(user *model.User) error
	Find(int64) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Edit(user *model.User) error
}