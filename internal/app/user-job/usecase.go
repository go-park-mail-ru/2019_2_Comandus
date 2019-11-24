package user_job

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateJob(user *model.User, job *model.Job) error
	FindJob(id int64) (*model.Job, error)
	GetAllJobs() ([]model.Job, error)
	EditJob(user *model.User, job *model.Job, id int64) error
	PatternSearch(string) ([]model.Job, error)
}
