package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func CreateCompanyOnServer(userId int64) (*company_grpc.Company, error){
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := company_grpc.NewCompanyHandlerClient(conn)
	comReq := &company_grpc.UserID{
		ID:                   userId,
	}
	company, err := client.CreateCompany(context.Background(),comReq)
	if err != nil {
		return nil, errors.Wrap(err, "client.CreateCompany")
	}
	return company, nil
}

func GetCompanyFromServer(id int64) (*company_grpc.Company, error) {
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}
	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := company_grpc.NewCompanyHandlerClient(conn)
	companyReq := &company_grpc.CompanyID{
		ID:		id,
	}

	currCompany, err := client.Find(context.Background(), companyReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currCompany, nil
}

func EditCompanyOnServer(id int64, company *model.Company)  error {
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := company_grpc.NewCompanyHandlerClient(conn)

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

	_, err = client.Edit(context.Background(), companyReq)
	if err != nil {
		return errors.Wrap(err, "userRep.Find()")
	}

	return nil
}