package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type UserClient struct {
	conn *grpc.ClientConn
}

func (c *UserClient) Connect() error {
	conn, err := grpc.Dial(":8087", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *UserClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *UserClient) GetUserFromServer(userID *user_grpc.UserID) (*user_grpc.User, error){
	client := user_grpc.NewUserHandlerClient(c.conn)

	req := &user_grpc.UserID{
		ID:                   userID.ID,
	}
	res, err := client.Find(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Find()")
	}

	return res, nil
}