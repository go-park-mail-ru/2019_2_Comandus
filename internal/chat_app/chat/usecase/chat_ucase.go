package chat_ucase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model/chat"
	"github.com/pkg/errors"
)

type ChatUsecase struct {
	rep chat.Repository
}

func NewChatUsecase(r chat.Repository) chat.Usecase {
	return &ChatUsecase{
		rep: r,
	}
}


func (u *ChatUsecase) CreateChat(newChat *model.Chat) (*model.Chat, error) {
	extChat, err := u.rep.FindCurrent(newChat.UserID, newChat.SupportID)
	if err == nil {
		return extChat, errors.New(chat.CHAT_CONFLICT_ERR)
	}

	if err := u.rep.Create(newChat); err != nil {
		return nil, errors.Wrapf(err, "chatRep.Create()")
	}
	return newChat, nil
}

func (u *ChatUsecase) Find(id int64) (*model.Chat, error) {
	c, err := u.rep.Find(id)
	if err != nil {
		return nil, errors.Wrapf(err, "chatRep.Find()")
	}
	return c, nil
}

func (u *ChatUsecase) List() ([]*model.Chat, error) {
	chats, err := u.rep.List(0,0)
	if err != nil {
		return nil, errors.Wrap(err, "chatRep.List()")
	}
	return chats, nil
}

func (u *ChatUsecase) ListByUser(id int64) ([]*model.Chat, error) {
	chats, err := u.rep.List(id,0)
	if err != nil {
		return nil, errors.Wrap(err, "chatRep.List()")
	}
	return chats, nil
}

func (u *ChatUsecase) ListByClient(id int64) ([]*model.Chat, error) {
	chats, err := u.rep.List(0,id)
	if err != nil {
		return nil, errors.Wrap(err, "chatRep.List()")
	}
	return chats, nil
}
