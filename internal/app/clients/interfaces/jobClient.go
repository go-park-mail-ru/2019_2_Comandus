package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc/job_grpc"
)

type ClientJob interface {
	GetJobFromServer(id int64) (*job_grpc.Job, error)
}
