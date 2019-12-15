package user_contract

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateContract(user *model.User, responseId int64, input *model.ContractInput) error
	SetAsDone(user *model.User, contractId int64) error
	ReviewContract(user *model.User, contractId int64, review *model.ReviewInput) error
	ReviewList(user *model.User) ([]model.Review, error)
	ContractList(user *model.User) ([]model.ContractOutput, error)
	Find(user *model.User, id int64) (*model.ContractOutput, error)
	ChangeStatus(user *model.User, id int64, status string) error
	TickWorkAsReady(user *model.User, id int64) error
}
