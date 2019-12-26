package chat

import (
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type Repository interface {
	Create(chat *model2.Chat) error
	Find(id int64) (*model2.Chat, error)
	FindCurrent(userId int64, clientId int64) (*model2.Chat, error)
	List(userId int64, clientId int64) ([]*model2.Chat, error)
	FindByProposal(id int64) (*model2.Chat, error)
}
