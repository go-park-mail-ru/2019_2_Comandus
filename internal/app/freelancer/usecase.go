package freelancer

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Create(int64) (*model.Freelancer, error)
	FindByUser(int64) (*model.FreelancerOutput, error)
	Find(int64) (*model.ExtendedOutputFreelancer, *model.FreelancerContracts, error)
	Edit(int64, *model.Freelancer) error
	PatternSearch(string, model.SearchParams) ([]model.ExtendFreelancer, error)
	FindPart(int, int) ([]model.ExtendFreelancer, error)
	FindNoLocation(int64) (*model.ExtendFreelancer, error)
	FindNoLocationByUser(id int64) (*model.Freelancer, error)
	GetRating(id int64) (int64, error)
}
