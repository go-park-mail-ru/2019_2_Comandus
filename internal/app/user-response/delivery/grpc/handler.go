package regrpc

import (
	"context"
	user_response "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/grpc/response_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
)

type ResponseServer struct {
	Ucase user_response.Usecase
}

func NewResponseServerGrpc(gserver *grpc.Server, ucase user_response.Usecase) {
	server := &ResponseServer{
		Ucase: ucase,
	}
	response_grpc.RegisterResponseHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *ResponseServer) TransformResponseRPC(response *model.Response) *response_grpc.Response {
	if response == nil {
		return nil
	}

	date := &timestamp.Timestamp{
		Seconds:              response.Date.Unix(),
		Nanos:                int32(response.Date.UnixNano()),
	}

	res := &response_grpc.Response{
		ID:                   response.ID,
		FreelancerId:         response.FreelancerId,
		JobId:                response.JobId,
		Files:                response.Files,
		Date:                 date,
		StatusManager:        response.StatusManager,
		StatusFreelancer:     response.StatusFreelancer,
		PaymentAmount:        response.PaymentAmount,
	}
	return res
}


func (s *ResponseServer) TransformResponseData(response *response_grpc.Response) *model.Response {
	// TODO: fix date
	res := &model.Response{
		ID:                   response.ID,
		FreelancerId:         response.FreelancerId,
		JobId:                response.JobId,
		Files:                response.Files,
		Date:                 time.Unix(response.Date.Seconds, int64(response.Date.Nanos)),
		StatusManager:        response.StatusManager,
		StatusFreelancer:     response.StatusFreelancer,
		PaymentAmount:        response.PaymentAmount,
	}
	return res
}

func (s *ResponseServer) Find(context context.Context, response *response_grpc.ResponseID) (*response_grpc.Response, error) {
	newResponse, err := s.Ucase.Find(response.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ucase.Find()")
	}
	res := s.TransformResponseRPC(newResponse)
	return res, nil
}

func (s *ResponseServer) Update(context context.Context, response *response_grpc.Response) (*response_grpc.Nothing, error) {
	resp := s.TransformResponseData(response)
	if err := s.Ucase.Update(resp); err != nil {
		return nil, err
	}
	return nil, nil
}