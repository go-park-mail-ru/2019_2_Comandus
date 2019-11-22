package cogrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CompanyServer struct {
	Ucase company.Usecase
}

func NewCompanyServerGrpc(gserver *grpc.Server, companyUcase company.Usecase) {
	companyServer := &CompanyServer{
		Ucase: companyUcase,
	}
	company_grpc.RegisterCompanyHandlerServer(gserver, companyServer)
	reflection.Register(gserver)
}

func (s *CompanyServer) TransformCompanyRPC(company *model.Company) *company_grpc.Company {
	if company == nil {
		return nil
	}

	res := &company_grpc.Company{
		ID:                   company.ID,
		CompanyName:          company.CompanyName,
		Site:                 company.Site,
		TagLine:              company.TagLine,
		Description:          company.Description,
		Country:              company.Country,
		City:                 company.City,
		Address:              company.Address,
		Phone:                company.Phone,
	}
	return res
}


func (s *CompanyServer) TransformCompanyData(company *company_grpc.Company) *model.Company {
	res := &model.Company{
		ID:          company.ID,
		CompanyName: company.CompanyName,
		Site:        company.Site,
		TagLine:     company.TagLine,
		Description: company.Description,
		Country:     company.Country,
		City:        company.City,
		Address:     company.Address,
		Phone:       company.Phone,
	}
	return res
}

func (s *CompanyServer) CreateCompany(context context.Context, userID *company_grpc.UserID) (*company_grpc.Company, error) {
	newCompany, err := s.Ucase.Create()
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.CreateUser")
	}

	res := s.TransformCompanyRPC(newCompany)
	return res, nil
}

func (s *CompanyServer) Find(context context.Context,company *company_grpc.CompanyID) (*company_grpc.Company, error) {
	newCompany, err := s.Ucase.Find(company.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ucase.Find()")
	}

	res := s.TransformCompanyRPC(newCompany)
	return res, nil
}
