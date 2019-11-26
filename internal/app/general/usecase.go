package general

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	SignUp(user *model.User) (int64, error)
	VerifyUser(user *model.User) (int64, error)
}

