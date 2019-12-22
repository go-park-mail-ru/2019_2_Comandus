package chat

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model/chat"

type Repository interface {
	Create(chat *model.Chat) error
	Find(id int64) (*model.Chat, error)
	FindCurrent(userId int64, clientId int64) (*model.Chat, error)
	List(userId int64, clientId int64) ([]*model.Chat, error)
}
