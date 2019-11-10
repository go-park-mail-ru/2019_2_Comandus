package freelancer

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	FindByUser(user *model.User) (*model.Freelancer, error)
	Find(id int64) (*model.Freelancer, error)
	Edit(user *model.User, freelancer *model.Freelancer) error
}