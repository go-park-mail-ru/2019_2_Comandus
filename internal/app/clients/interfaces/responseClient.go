package clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/proposal/delivery/grpc/response_grpc"
)

type ClientResponse interface {
	GetResponseFromServer(id int64) (*response_grpc.Response, error)
	UpdateResponseOnServer(response *response_grpc.Response) error
}
