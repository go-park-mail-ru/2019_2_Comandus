package manager

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Usecase interface {
	Create(int64, int64) (*model.HireManager, error)
}
