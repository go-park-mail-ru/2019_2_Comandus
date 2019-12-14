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

func (s *CompanyServer) TransformCompanyOutputRPC(company *model.CompanyOutput) *company_grpc.CompanyOutput {
	if company == nil {
		return nil
	}

	res := &company_grpc.CompanyOutput{
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
		return nil, errors.Wrap(err, "Ucase.CreateUser")
	}

	res := s.TransformCompanyRPC(newCompany)
	return res, nil
}

func (s *CompanyServer) Find(context context.Context, company *company_grpc.CompanyID) (*company_grpc.CompanyOutput, error) {
	newCompany, err := s.Ucase.Find(company.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ucase.Find()")
	}

	res := s.TransformCompanyOutputRPC(newCompany)
	return res, nil
}

func (s *CompanyServer) Edit(context context.Context, inputCompany *company_grpc.CompanyWithUser) (*company_grpc.Nothing, error) {
	myCompany := s.TransformCompanyData(inputCompany.MyCompany)
	if _, err := s.Ucase.Edit(inputCompany.UserID, myCompany); err != nil {
		return nil, errors.Wrap(err, "Ucase.Edit()")
	}
	return nil, nil
}
