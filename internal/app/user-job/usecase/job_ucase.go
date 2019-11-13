package jobUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"time"
)

type JobUsecase struct {
	managerRep		manager.Repository
	jobRep			user_job.Repository
}

func NewJobUsecase(m manager.Repository, j user_job.Repository) user_job.Usecase {
	return &JobUsecase{
		managerRep:		m,
		jobRep:			j,
	}
}

func (u *JobUsecase) CreateJob(currUser * model.User, job *model.Job) error {
	if !currUser.IsManager() {
		return errors.New("current user is not a manager: ")
	}

	currManager, err := u.managerRep.FindByUser(currUser.ID)
	if err != nil {
		return errors.Wrapf(err, "managerRep.FindByUser(): ")
	}

	job.Date = time.Now()
	if err = u.jobRep.Create(job, currManager); err != nil {
		return errors.Wrapf(err, "jobRep.Create(): ")
	}
	return nil
}

func (u *JobUsecase) FindJob(id int64) (*model.Job, error) {
	job, err := u.jobRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "jobRep.Find(): ")
	}
	return job, nil
}

func (u *JobUsecase) GetAllJobs() ([]model.Job, error) {
	jobs, err := u.jobRep.List()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		return nil, errors.Wrap(err, "jobRep.List(): ")
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

	currManager, err := u.managerRep.FindByUser(user.ID)
	if err != nil {
		return errors.Wrap(err, "managerRep.FindByUser(): ")
	}

	if job.HireManagerId != currManager.ID {
		return errors.New("no access for current manager")
	}

	inputJob.ID = job.ID
	inputJob.HireManagerId = job.HireManagerId

	if err := u.jobRep.Edit(inputJob); err != nil {
		return errors.Wrapf(err, "HandleEditProfile<-JobEdit")
	}
	return nil
}