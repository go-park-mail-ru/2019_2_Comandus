package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

type CompanyClient struct {
	conn *grpc.ClientConn
}

func (c *CompanyClient) Connect() error {
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *CompanyClient) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		log.Println("conn.Close()", err)
	}
	return nil
}

func (c *CompanyClient) CreateCompanyOnServer(userId int64) (*company_grpc.Company, error) {
	client := company_grpc.NewCompanyHandlerClient(c.conn)
	comReq := &company_grpc.UserID{
		ID: userId,
	}
	company, err := client.CreateCompany(context.Background(), comReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateCompany")
	}
	return company, nil
}

func (c *CompanyClient) GetCompanyFromServer(id int64) (*company_grpc.Company, error) {
	client := company_grpc.NewCompanyHandlerClient(c.conn)
	companyReq := &company_grpc.CompanyID{
		ID: id,
	}

	currCompany, err := client.Find(context.Background(), companyReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currCompany, nil
}

func (c *CompanyClient) EditCompanyOnServer(id int64, company *model.Company)  error {
	client := company_grpc.NewCompanyHandlerClient(c.conn)

	grpccompany := &company_grpc.Company{
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

	companyReq := &company_grpc.CompanyWithUser{
		MyCompany:            grpccompany,
		UserID:               id,
	}

	_, err := client.Edit(context.Background(), companyReq)
	if err != nil {
		return errors.Wrap(err, "userRep.Find()")
	}

	return nil
}
