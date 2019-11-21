package jgrpc

import (
	"context"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc/job_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
)

type JobServer struct {
	Ucase user_job.Usecase
}

func NewJobServerGrpc(gserver *grpc.Server, jobUcase user_job.Usecase) {
	server := &JobServer{
		Ucase: jobUcase,
	}
	job_grpc.RegisterJobHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *JobServer) TransformJobRPC(job *model.Job) *job_grpc.Job {
	if job == nil {
		return nil
	}

	date := &timestamp.Timestamp{
		Seconds:              job.Date.Unix(),
		Nanos:                int32(job.Date.UnixNano()),
	}

	res := &job_grpc.Job{
		ID:                   job.ID,
		HireManagerId:        job.HireManagerId,
		Title:                job.Title,
		Description:          job.Description,
		Files:                job.Files,
		SpecialityId:         job.SpecialityId,
		ExperienceLevelId:    job.ExperienceLevelId,
		PaymentAmount:        job.PaymentAmount,
		Country:              job.Country,
		City:                 job.City,
		JobTypeId:            job.JobTypeId,
		Date:                 date,
		Status:               job.Status,
	}
	return res
}


func (s *JobServer) TransformJobData(job *job_grpc.Job) *model.Job {
	// TODO: fix date
	date := time.Time{}

	res := &model.Job{
		ID:                job.ID,
		HireManagerId:     job.HireManagerId,
		Title:             job.Title,
		Description:       job.Description,
		Files:             job.Files,
		SpecialityId:      job.SpecialityId,
		ExperienceLevelId: job.ExperienceLevelId,
		PaymentAmount:     job.PaymentAmount,
		Country:           job.Country,
		City:              job.City,
		JobTypeId:         job.JobTypeId,
		Date:              date,
		Status:            job.Status,
	}
	return res
}

func (s *JobServer) Find(context context.Context, jobId *job_grpc.JobID) (*job_grpc.Job, error) {
	currJob, err := s.Ucase.FindJob(jobId.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ucase.FindByUser()")
	}
	res := s.TransformJobRPC(currJob)
	return res, nil
}