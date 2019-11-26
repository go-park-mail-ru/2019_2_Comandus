package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/auth/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func CreateUserOnServer(data *model.User) (*auth_grpc.User, error) {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := auth_grpc.NewAuthHandlerClient(conn)
	userReq := &auth_grpc.User{
		Email:              data.Email,
		Password:			data.Password,
		FirstName:			data.FirstName,
		SecondName:			data.SecondName,
		UserType:			data.UserType,
	}

	user, err := client.CreateUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateUser")
	}

	return user, nil
}

func VerifyUserOnServer(user *model.User) (int64, error){
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return 0, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := auth_grpc.NewAuthHandlerClient(conn)

	mes, err := client.VerifyUser(context.Background(), &auth_grpc.UserRequest{
		Email:                user.Email,
		Password:             user.Password,
	})

	if err != nil {
		return 0, errors.Wrap(err, "client.VerifyUser()")
	}

	return mes.ID, nil
}
