package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat/delivery/grpc/chat_grpc"
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type ChatClient struct {
	conn *grpc.ClientConn
}

func (c *ChatClient) Connect() error {
	conn, err := grpc.Dial(CHAT_PORT, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *ChatClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *ChatClient) CreateChatOnServer(chat *model2.Chat) error {
	client := chat_grpc.NewChatHandlerClient(c.conn)

	grpcchat := &chat_grpc.Chat{
		ID:                   chat.ID,
		User:                 chat.Freelancer,
		Support:              chat.Manager,
		Name:                 chat.Name,
		Proposal:             chat.ProposalId,
	}

	if _, err := client.Create(context.Background(), grpcchat); err != nil {
		return errors.Wrap(err, "client.CreateChat()")
	}
	return nil
}
