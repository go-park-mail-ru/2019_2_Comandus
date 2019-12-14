package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc/auth_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

const (
	AUTH_PORT       = ":8081"
	COMPANY_PORT    = ":8082"
	FREELANCER_PORT = ":8083"
	JOB_PORT        = ":8084"
	LOCATION_PORT   = ":8085"
	MANAGER_PORT    = ":8086"
	RESPONSE_PORT   = ":8087"
	USER_PORT       = ":8088"
)

type AuthClient struct {
	conn *grpc.ClientConn
}

func (c *AuthClient) Connect() error {
	conn, err := grpc.Dial(AUTH_PORT, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *AuthClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *AuthClient) CreateUserOnServer(data *model.User) (*auth_grpc.User, error) {
	client := auth_grpc.NewAuthHandlerClient(c.conn)
	userReq := &auth_grpc.User{
		Email:      data.Email,
		Password:   data.Password,
		FirstName:  data.FirstName,
		SecondName: data.SecondName,
		UserType:   data.UserType,
	}

	user, err := client.CreateUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateUser")
	}

	return user, nil
}

func (c *AuthClient) VerifyUserOnServer(user *model.User) (int64, error) {
	client := auth_grpc.NewAuthHandlerClient(c.conn)

	mes, err := client.VerifyUser(context.Background(), &auth_grpc.UserRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		return 0, errors.Wrap(err, "client.VerifyUser()")
	}

	return mes.ID, nil
}
