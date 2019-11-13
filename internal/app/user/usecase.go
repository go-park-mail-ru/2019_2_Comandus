package user

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateUser(*model.User) error
	VerifyUser(*model.User) (int64, error)
	EditUser(newUser *model.User, oldUser *model.User) error
	EditUserPassword(passwords *model.BodyPassword, user *model.User) error
	Find(int64) (*model.User, error)
	SetUserType(user *model.User, userType string) error
	GetRoles(user *model.User) ([]*model.Role, error)
	GetAvatar(user *model.User) ([]byte, error)
}
