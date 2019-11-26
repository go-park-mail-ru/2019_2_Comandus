package user_job

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(j *model.Job) error
	Find(int64) (*model.Job, error)
	Edit(job *model.Job) error
	List() ([]model.Job, error)
	ListOnPattern (string) ([]model.Job, error)
}
