package chat

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

const (
	CHAT_CONFLICT_ERR = "chat_conflict"
)


type Usecase interface {
	CreateChat(newChat *model.Chat) (*model.Chat, error)
	Find(id int64) (*model.Chat, error)
	FindByProposal(id int64) (*model.Chat, error)
	List() ([]*model.Chat, error)
}