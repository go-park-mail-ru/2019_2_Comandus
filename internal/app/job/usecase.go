package job

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	CreateJob(user *model.User, job *model.Job) error
	FindJob(id int64) (*model.Job, error)
	GetAllJobs() ([]model.Job, error)
	GetMyJobs(int64) ([]model.Job, error)
	EditJob(user *model.User, job *model.Job, id int64) error
	MarkAsDeleted(id int64, user *model.User) error
	PatternSearch(string, model.SearchParams) ([]model.Job, error)
	GetUserIDByJobID(jobID int64) (int64, error)
	ChangeStatus(jobID int64, status string, userID int64) error
}
