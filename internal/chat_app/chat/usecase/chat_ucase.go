package chat_ucase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
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
	extChat, err := u.rep.FindCurrent(newChat.Freelancer, newChat.Manager)
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

func (u *ChatUsecase) FindByProposal(id int64) (*model.Chat, error) {
	c, err := u.rep.FindByProposal(id)
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

