package general_ucase

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc/user_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type GeneralUsecase struct {
}

func NewGeneralUsecase() general.Usecase {
	return &GeneralUsecase{
	}
}

func (u *GeneralUsecase) SignUp(data *model.User) error {
	var conns [4]*grpc.ClientConn
	var err error
	for i := 0; i < 4; i++ {
		port := ":808" + strconv.Itoa(i+1)
		conns[i], err = grpc.Dial(port, grpc.WithInsecure())
		if err != nil {
			return errors.Wrap(err, "grpc.Dial()")
		}
	}

	defer func(){
		for i := 0; i < 4; i++ {
			if err := conns[i].Close(); err != nil {
				// TODO: use zap logger
				log.Println("GeneralUsecase<-CreateUser:", err)
			}
		}
	}()

	client := user_grpc.NewUserHandlerClient(conns[0])
	user, err := client.CreateUser(context.Background(), &user_grpc.UserRequest{
		Email:                data.Email,
		Password:             data.Password,
	})
	if err != nil {
		return errors.Wrap(err, "client.CreateUser")
	}

	client1 := company_grpc.NewCompanyHandlerClient(conns[1])
	company, err := client1.CreateCompany(context.Background(), &company_grpc.UserID{
		ID:                   user.ID,
	})
	if err != nil {
		return errors.Wrap(err, "client.CreateCompany")
	}

	client2 := freelancer_grpc.NewFreelancerHandlerClient(conns[2])
	_,err = client2.CreateFreelancer(context.Background(), &freelancer_grpc.UserID{
		ID:                   user.ID,
	})

	client3 := manager_grpc.NewManagerHandlerClient(conns[3])
	_, err = client3.CreateManager(context.Background(), &manager_grpc.Info{
		UserID:user.ID,
		CompanyID:company.ID,
	})

	return nil
}

func (u *GeneralUsecase) VerifyUser(user *model.User) (int64, error) {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return 0, errors.Wrap(err, "grpc.Dial()")
	}

	client := user_grpc.NewUserHandlerClient(conn)

	mes, err := client.VerifyUser(context.Background(), &user_grpc.UserRequest{
		Email:                user.Email,
		Password:             user.Password,
	})

	if err != nil {
		return 0, errors.Wrap(err, "client.VerifyUser")
	}

	return mes.ID, nil
}
