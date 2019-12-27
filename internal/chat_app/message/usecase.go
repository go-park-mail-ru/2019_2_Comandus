package message

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type Usecase interface {
	Create(message *model.Message) error
	List(chatId int64) ([]*model.Message, error)
	UpdateStatus(chatId int64, userId int64) error
}
