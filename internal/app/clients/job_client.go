package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc/job_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type JobClient struct {
	conn *grpc.ClientConn
}

func (c *JobClient) Connect() error {
	conn, err := grpc.Dial(JOB_PORT, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *JobClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *JobClient) GetJobFromServer(id int64) (*job_grpc.Job, error) {
	client := job_grpc.NewJobHandlerClient(c.conn)
	jobReq := &job_grpc.JobID{
		ID:		id,
	}

	currJob, err := client.Find(context.Background(), jobReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currJob, nil
}