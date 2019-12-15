package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
)

type ClientFreelancer interface {
	CreateFreelancerOnServer(userId int64) (*freelancer_grpc.Freelancer, error)
	GetFreelancerByUserFromServer(id int64) (*freelancer_grpc.Freelancer, error)
	GetFreelancerFromServer(id int64) (*freelancer_grpc.ExtendedFreelancer, error)
}
