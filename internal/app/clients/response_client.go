package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/grpc/response_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func GetResponseFromServer(id int64) (*response_grpc.Response, error){
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := response_grpc.NewResponseHandlerClient(conn)

	req := &response_grpc.ResponseID {
		ID:		id,
	}

	currResponse, err := client.Find(context.Background(), req)
	if err != nil {
		return nil, errors.Wrapf(err, "client.Find()")
	}
	return currResponse, nil
}
