package freelancer

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(freelancer *model.Freelancer) error
	Find(int64) (*model.Freelancer, error)
	FindByUser(int64) (*model.Freelancer, error)
	Edit(freelancer *model.Freelancer) error
	ListOnPattern (string) ([]model.Freelancer, error)
}

