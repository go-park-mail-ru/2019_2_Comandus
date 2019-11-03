package store

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	//"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
)

type UserRepository interface {
	Create(user *model.User) error
	Find(int64) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Edit(user *model.User) error
}

type FreelancerRepository interface {
	Create(freelancer *model.Freelancer) error
	Find(int64) (*model.Freelancer, error)
	Edit(freelancer *model.Freelancer) error
}

type ManagerRepository interface {
	Create(manager *model.HireManager) error
	Find(int64) (*model.Freelancer, error)
	Edit(manager *model.HireManager) error
}

type JobRepository interface {
	Create(job *model.Job) error
	Find(int64) (*model.Job, error)
	Edit(job *model.Job) error
	List() ([]model.Job, error)
}

type ResponseRepository interface {
	Create(response *model.Response) error
	Edit(response *model.Response) error
	ListForFreelancer(int64) ([]model.Response, error)
	ListForManager(int64) ([]model.Response, error)
	Find(int64) (*model.Response, error)
}

type CompanyRepository interface {
	Create(company *model.Company) error
	Find(int64) (*model.Company, error)
}