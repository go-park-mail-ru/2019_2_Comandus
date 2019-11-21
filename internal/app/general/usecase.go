package general

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	SignUp(user *model.User) error
	VerifyUser(user *model.User) (int64, error)
}

