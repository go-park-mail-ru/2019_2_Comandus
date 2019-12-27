package clients

import (
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type ChatClient interface {
	CreateChatOnServer(chat *model2.Chat) error
}
