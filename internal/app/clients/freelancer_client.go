package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type FreelancerClient struct{
	conn *grpc.ClientConn
}

func (c *FreelancerClient) Connect() error {
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *FreelancerClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *FreelancerClient) CreateFreelancerOnServer(userId int64) (*freelancer_grpc.Freelancer, error) {
	client := freelancer_grpc.NewFreelancerHandlerClient(c.conn)
	fReq := &freelancer_grpc.UserID{
		ID: userId,
	}
	freelancer, err := client.CreateFreelancer(context.Background(), fReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateFreelancer()")
	}
	return freelancer, nil
}

func (c *FreelancerClient) GetFreelancerByUserFromServer(id int64) (*freelancer_grpc.Freelancer, error) {
	client := freelancer_grpc.NewFreelancerHandlerClient(c.conn)
	userReq := &freelancer_grpc.UserID{
		ID: id,
	}

	currFreelancer, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currFreelancer, nil
}

func (c *FreelancerClient) GetFreelancerFromServer(id int64) (*freelancer_grpc.Freelancer, error) {
	client := freelancer_grpc.NewFreelancerHandlerClient(c.conn)
	req := &freelancer_grpc.FreelancerID{
		ID:		id,
	}

	currFreelancer, err := client.Find(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currFreelancer, nil
}