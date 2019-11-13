package manager

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(manager *model.HireManager) error
	Find(int64) (*model.HireManager, error)
	FindByUser(int64) (*model.HireManager, error)
	Edit(manager *model.HireManager) error
	GetCompanyIDByUserID(accountId int64) (int64, error)
}
