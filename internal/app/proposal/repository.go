package proposal

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(response *model.Response) error
	Edit(response *model.Response) error
	ListForFreelancer(int64) ([]model.ExtendResponse, error)
	ListForManager(int64) ([]model.ExtendResponse, error)
	Find(int64) (*model.Response, error)
	ListResponsesOnJobID(jobID int64) ([]model.ExtendResponse, error)
	CheckForHavingResponse(jobID int64, freelID int64) (bool, error)
}
