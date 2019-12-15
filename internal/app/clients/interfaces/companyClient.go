package clients

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"

type CompanyClient interface {
	CreateCompanyOnServer(int64) (*company_grpc.Company, error)
	GetCompanyFromServer(int64) (*company_grpc.CompanyOutput, error)
}
