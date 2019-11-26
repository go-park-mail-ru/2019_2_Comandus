package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)


func GetUserFromServer(userID *user_grpc.UserID) (*user_grpc.User, error){
	conn, err := grpc.Dial(":8087", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close", err)
		}
	}()

	client := user_grpc.NewUserHandlerClient(conn)

	req := &user_grpc.UserID{
		ID:                   userID.ID,
	}
	res, err := client.Find(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Find()")
	}

	return res, nil
}