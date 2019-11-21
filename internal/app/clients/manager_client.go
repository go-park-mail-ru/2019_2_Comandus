package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func CreateManagerOnServer(userId int64, companyId int64) (*manager_grpc.Manager, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := manager_grpc.NewManagerHandlerClient(conn)
	manReq := &manager_grpc.Info{
		UserID:		userId,
		CompanyID:	companyId,
	}
	manager, err := client.CreateManager(context.Background(), manReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateManager()")
	}
	return manager, nil
}

func GetManagerByUserFromServer(id int64) (*manager_grpc.Manager, error){
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := manager_grpc.NewManagerHandlerClient(conn)

	userReq := &manager_grpc.UserID{
		ID:		id,
	}

	currManager, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrapf(err, "client.FindByUser()")
	}
	return currManager, nil
}

func GetManagerFromServer(id int64) (*manager_grpc.Manager, error){
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := manager_grpc.NewManagerHandlerClient(conn)

	userReq := &manager_grpc.ManagerID {
		ID:		id,
	}

	currManager, err := client.Find(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrapf(err, "client.Find()")
	}
	return currManager, nil
}
