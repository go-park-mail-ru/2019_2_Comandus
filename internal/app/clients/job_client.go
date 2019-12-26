package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/job/delivery/grpc/job_grpc"
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
		ID: id,
	}

	currJob, err := client.Find(context.Background(), jobReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currJob, nil
}

func (c *JobClient) GetUserIDByJobID(jobid int64) (int64, error) {
	client := job_grpc.NewJobHandlerClient(c.conn)
	jobReq := &job_grpc.JobID{
		ID: jobid,
	}

	uID, err := client.GetUserIDFromJobID(context.Background(), jobReq)
	if err != nil {
		return -1 , errors.Wrap(err, "userRep.Find()")
	}

	return uID.ID, nil
}

func (c *JobClient) GetTags() ([]string, error) {
	client := job_grpc.NewJobHandlerClient(c.conn)

	nothing := &job_grpc.Nothing{}
	jobs, err := client.GetTags(context.Background(), nothing)
	if err != nil {
		return nil, err
	}
	return jobs.Tags, err
}
