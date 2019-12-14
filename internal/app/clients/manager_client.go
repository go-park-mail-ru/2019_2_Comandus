package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type ManagerClient struct {
	conn *grpc.ClientConn
}

func (c *ManagerClient) Connect() error {
	conn, err := grpc.Dial(MANAGER_PORT, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *ManagerClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *ManagerClient) CreateManagerOnServer(userId int64, companyId int64) (*manager_grpc.Manager, error) {
	client := manager_grpc.NewManagerHandlerClient(c.conn)
	manReq := &manager_grpc.Info{
		UserID:    userId,
		CompanyID: companyId,
	}
	manager, err := client.CreateManager(context.Background(), manReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateManager()")
	}
	return manager, nil
}

func (c *ManagerClient) GetManagerByUserFromServer(id int64) (*manager_grpc.Manager, error) {
	client := manager_grpc.NewManagerHandlerClient(c.conn)

	userReq := &manager_grpc.UserID{
		ID: id,
	}

	currManager, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrapf(err, "client.FindByUser()")
	}
	return currManager, nil
}

func (c *ManagerClient) GetManagerFromServer(id int64) (*manager_grpc.Manager, error) {
	client := manager_grpc.NewManagerHandlerClient(c.conn)

	userReq := &manager_grpc.ManagerID{
		ID: id,
	}

	currManager, err := client.Find(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrapf(err, "client.Find()")
	}
	return currManager, nil
}
