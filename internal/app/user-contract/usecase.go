package user_contract

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateContract(user *model.User, responseId int64) error
	SetAsDone(user *model.User, contractId int64) error
	ReviewContract(user *model.User, contractId int64, grade int) error
}
