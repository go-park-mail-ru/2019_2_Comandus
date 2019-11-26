package clients

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"

type ClientManager interface {
CreateManagerOnServer(userId int64, companyId int64) (*manager_grpc.Manager, error)
GetManagerByUserFromServer(id int64) (*manager_grpc.Manager, error)
GetManagerFromServer(id int64) (*manager_grpc.Manager, error)
}
