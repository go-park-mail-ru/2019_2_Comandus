package store

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	//"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
)

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type FreelancerRepository interface {
	Create(freelancer *model.Freelancer) error
	Find(int) (*model.Freelancer, error)
}

type ManagerRepository interface {
	Create(manager *model.HireManager) error
	Find(int) (*model.Freelancer, error)
}
