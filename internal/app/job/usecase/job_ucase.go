package jobUcase

import (
	clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type JobUsecase struct {
	jobRep        user_job.Repository
	managerClient clients.ManagerClient
}

func NewJobUsecase(j user_job.Repository, mClient clients.ManagerClient) user_job.Usecase {
	return &JobUsecase{
		jobRep:        j,
		managerClient: mClient,
	}
}

func (u *JobUsecase) CreateJob(currUser *model.User, job *model.Job) error {
	if !currUser.IsManager() {
		return errors.New("current user is not a manager")
	}

	currManager, err := u.managerClient.GetManagerByUserFromServer(currUser.ID)
	if err != nil {
		return errors.Wrapf(err, "getManagerByUserFromServer()")
	}

	job.HireManagerId = currManager.ID
	job.Date = time.Now()
	job.Status = model.JobStateOpened
	if err = u.jobRep.Create(job); err != nil {
		return errors.Wrapf(err, "jobRep.Create()")
	}
	return nil
}

func (u *JobUsecase) FindJob(id int64) (*model.Job, error) {
	job, err := u.jobRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "jobRep.Find()")
	}
	return job, nil
}

func (u *JobUsecase) GetAllJobs() ([]model.Job, error) {
	jobs, err := u.jobRep.List()
	if err != nil {
		return nil, errors.Wrap(err, "jobRep.List()")
	}
	return jobs, nil
}

func (u *JobUsecase) EditJob(user *model.User, inputJob *model.Job, id int64) error {
	if !user.IsManager() {
		return errors.New("only manager can edit job")
	}

	job, err := u.jobRep.Find(id)
	if err != nil {
		return errors.Wrapf(err, "jobRep.Find(): ")
	}

	currManager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "getManagerByUserFromServer()")
	}

	if job.HireManagerId != currManager.ID {
		return errors.New("no access for current manager")
	}

	inputJob.ID = job.ID
	inputJob.HireManagerId = job.HireManagerId

	if err := u.jobRep.Edit(inputJob); err != nil {
		return errors.Wrapf(err, "jobRep.Edit()")
	}
	return nil
}

func (u *JobUsecase) MarkAsDeleted(id int64, user *model.User) error {
	job, err := u.jobRep.Find(id)
	if err != nil {
		return errors.Wrap(err, "jobRep.Find()")
	}

	if !user.IsManager() {
		return errors.New("only manager can delete job")
	}

	manager, err := u.managerClient.GetManagerByUserFromServer(user.ID)
	if err != nil {
		return errors.Wrap(err, "clients.GetManagerByUserFromServer()")
	}

	if job.HireManagerId != manager.ID {
		return errors.Wrap(err, "no access for current manager")
	}

	job.Status = model.JobStateDeleted
	if err := u.jobRep.Edit(job); err != nil {
		return errors.Wrap(err, "jobRep.Edit()")
	}

	return nil
}

func (u *JobUsecase) PatternSearch(pattern string, params model.SearchParams) ([]model.Job, error) {
	if (params.MaxPaymentAmount < params.MinPaymentAmount) ||
		(params.Limit < 0) {
		return nil, errors.New("wrong input parameters")
	}

	jobs, err := u.jobRep.ListOnPattern(pattern, params)
	if err != nil {
		return nil, errors.Wrap(err, "PatternSearch()")
	}
	return jobs, nil
}

func (u *JobUsecase) GetMyJobs(managerID int64) ([]model.Job, error) {
	jobs, err := u.jobRep.ListMyJobs(managerID)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetMyJobs<-ListMyJobs")
		return nil, err
	}
	return jobs, nil
}

func (u *JobUsecase) GetUserIDByJobID(jobID int64) (int64, error) {
	return u.jobRep.GetUserIDByJobID(jobID)
}

func (u *JobUsecase) ChangeStatus(jobID int64, status string, userID int64) error{
	job, err := u.jobRep.Find(jobID)
	if err != nil {
		return errors.Wrap(err, "jobRep.Find()")
	}

	manager, err := u.managerClient.GetManagerByUserFromServer(userID)
	if err != nil {
		return errors.Wrap(err, "clients.GetManagerByUserFromServer()")
	}

	if job.HireManagerId != manager.ID {
		return errors.Wrap(err, "no access for current manager")
	}

	if err := u.jobRep.ChangeStatus(jobID, status); err != nil {
		return errors.Wrap(err, "jobRep.ChangeStatus()")
	}

	return nil
}