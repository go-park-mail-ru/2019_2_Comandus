package message

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model/chat"

type Usecase interface {
	Create(message *model.Message) error
	ListByUser(userId int64, chatId int64) ([]*model.Message, error)
	ListBySupport(supportId int64, chatId int64) ([]*model.Message, error)
}
