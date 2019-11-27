package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/grpc/response_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)


type ResponseClient struct {
	conn *grpc.ClientConn
}

func (c *ResponseClient) Connect() error {
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *ResponseClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *ResponseClient) GetResponseFromServer(id int64) (*response_grpc.Response, error){
	client := response_grpc.NewResponseHandlerClient(c.conn)

	req := &response_grpc.ResponseID {
		ID:		id,
	}

	currResponse, err := client.Find(context.Background(), req)
	if err != nil {
		return nil, errors.Wrapf(err, "client.Find()")
	}
	return currResponse, nil
}
