package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func CreateFreelancerOnServer(userId int64) (*freelancer_grpc.Freelancer, error){
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := freelancer_grpc.NewFreelancerHandlerClient(conn)
	fReq := &freelancer_grpc.UserID{
		ID:		userId,
	}
	freelancer, err := client.CreateFreelancer(context.Background(), fReq)
	return freelancer, nil
}

func GetFreelancerByUserFromServer(id int64) (*freelancer_grpc.Freelancer, error) {
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := freelancer_grpc.NewFreelancerHandlerClient(conn)
	userReq := &freelancer_grpc.UserID{
		ID:		id,
	}

	currFreelancer, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currFreelancer, nil
}