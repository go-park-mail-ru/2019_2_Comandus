package proposal

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateResponse(user *model.User, response *model.Response, jobId int64) error
	GetResponses(user *model.User) ([]model.ExtendResponse, error)
	AcceptResponse(user *model.User, responseId int64) error
	DenyResponse(user *model.User, responseId int64) error
	CancelResponse(user *model.User, responseId int64) error
	Find(id int64) (*model.Response, error)
	GetResponse(id int64) (*model.ResponseOutputWithFreel, error)
	GetResponsesOnJobID(jobID int64) ([]model.ExtendResponse, error)
	Update(response *model.Response) error
}
