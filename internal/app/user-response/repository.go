package user_response

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(response *model.Response) error
	Edit(response *model.Response) error
	ListForFreelancer(int64) ([]model.Response, error)
	ListForManager(int64) ([]model.Response, error)
	Find(int64) (*model.Response, error)
}
