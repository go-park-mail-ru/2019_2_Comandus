package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc/job_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func GetJobFromServer(id int64) (*job_grpc.Job, error) {
	conn, err := grpc.Dial(":8085", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := job_grpc.NewJobHandlerClient(conn)
	jobReq := &job_grpc.JobID{
		ID:		id,
	}

	currJob, err := client.Find(context.Background(), jobReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currJob, nil
}