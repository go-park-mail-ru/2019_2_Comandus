package freelancer

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Create(int64) (*model.Freelancer, error)
	FindByUser(int64) (*model.Freelancer, error)
	Find(int64) (*model.Freelancer, error)
	Edit(*model.Freelancer, *model.Freelancer) error
	FindWithLocation(int64) (*model.FreelancerOutput, error)
	FindByUserWithLocation(int64) (*model.FreelancerOutput, error)
}