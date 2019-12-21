package contract

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

type Repository interface {
	Create(contract *model.Contract) error
	Edit(contract *model.Contract) error
	List(int64, string) ([]model.Contract, error)
	Find(int64) (*model.Contract, error)
	ChangeStatus(contractID int64, newStatus string) error
	ChangeStatusWorkAsReady(conractID int64) error
	GetClosedContracts(id int64) ([]model.PublicContractVersion, error)
}
