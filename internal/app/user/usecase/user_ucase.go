package userUcase

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc/company_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc/freelancer_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"os"
)

type UserUsecase struct {
	userRep			user.Repository
}

func NewUserUsecase(u user.Repository) user.Usecase {
	return &UserUsecase{
		userRep:		u,
	}
}

func (u *UserUsecase) CreateUser(data *model.User) error {
	if err := data.Validate(); err != nil {
		return errors.Wrap(err, "user.Validate()")
	}

	if err := data.BeforeCreate(); err != nil {
		return errors.Wrap(err, "user.BeforeCreate()")
	}

	if err := u.userRep.Create(data); err != nil {
		return errors.Wrap(err, "userRep.Create()")
	}

	return nil
}

func (u * UserUsecase) EditUser(new *model.User, old * model.User) error {
	new.ID = old.ID

	if old.Email != new.Email {
		return errors.Wrap(errors.New("can't change email"), "EditUser")
	}

	new.UserType = old.UserType
	new.EncryptPassword = old.EncryptPassword

	// TODO: можно аватар через отдельный метод изменять// нужно ли это тут вообще если аватар грузится в другом месте
	new.Avatar = old.Avatar
	if err := u.userRep.Edit(new); err != nil {
		return errors.Wrap(err, "userRep.Edit()")
	}
	return nil
}

func (u *UserUsecase) EditUserPassword(passwords *model.BodyPassword, user *model.User) error {
	if passwords.NewPassword != passwords.NewPasswordConfirmation {
		return errors.New("new passwords are different")
	}

	if !user.ComparePassword(passwords.Password) {
		err := errors.New("wrong old password")
		return errors.Wrapf(err, "model.user.ComparePassword")
	}

	newEncryptPassword, err := model.EncryptString(passwords.NewPassword)
	if err != nil {
		return errors.Wrap(err,"model.EncryptString")
	}
	user.EncryptPassword = newEncryptPassword

	if err := u.userRep.Edit(user); err != nil {
		return errors.Wrapf(err, "userRep.Edit")
	}
	return nil
}

func (u *UserUsecase) GetAvatar(user *model.User) ([]byte, error) {
	if user.Avatar != nil {
		return user.Avatar, nil
	}

	var openFile *os.File

	// TODO: create default user in database, get default image from it
	filename := "../../../../store/avatars/default.png"
	openFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Open")
	}

	defer func() {
		if err := openFile.Close(); err != nil {
			// TODO: write in correct logger
			log.Println(errors.Wrap(err, "GetAvatar<-Close()"))
		}
	}()

	avatar := make([]byte, 0)
	if _, err := openFile.Read(avatar); err != nil {
		return nil, errors.Wrap(err, "Read()")
	}

	return avatar, nil
}

func (u *UserUsecase) getFreelancerByUserFromServer(id int64) (*freelancer_grpc.Freelancer, error) {
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := freelancer_grpc.NewFreelancerHandlerClient(conn)
	userReq := &freelancer_grpc.UserID{
		ID:		id,
	}

	currFreelancer, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currFreelancer, nil
}

func (u *UserUsecase) getManagerByUserFromServer(id int64) (*manager_grpc.Manager, error) {
	conn, err := grpc.Dial(":8084", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}
	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := manager_grpc.NewManagerHandlerClient(conn)

	userReq := &manager_grpc.UserID{
		ID:		id,
	}

	currManager, err := client.FindByUser(context.Background(), userReq)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	return currManager, nil
}

func (u *UserUsecase) getCompanyFromServer(id int64) (*company_grpc.Company, error) {
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

func (u *UserUsecase) Find(id int64) (*model.User, error) {
	user, err := u.userRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	currFreelancer, err := u.getFreelancerByUserFromServer(id)
	if err != nil {
		return nil, errors.Wrap(err, "getFreelancerByUserFromServer()")
	}

	currManager, err := u.getManagerByUserFromServer(id)
	if err != nil {
		return nil, errors.Wrap(err, "getManagerByUserFromServer()")
	}

	user.FreelancerId = currFreelancer.ID
	user.HireManagerId = currManager.ID
	user.CompanyId = currManager.CompanyId

	return user, nil
}

func (u *UserUsecase) SetUserType(user *model.User, userType string) error {
	if err := user.SetUserType(userType); err != nil {
		return errors.Wrap(err, "SetUserType()")
	}

	if err := u.userRep.Edit(user); err != nil {
		return errors.Wrap(err, "userRep.Edit()")
	}
	return nil
}

func (u *UserUsecase) VerifyUser(currUser *model.User) (int64, error) {
	us, err := u.userRep.FindByEmail(currUser.Email)
	if err != nil {
		return 0, errors.Wrapf(err, "userRep.FindByEmail()")
	}

	if !us.ComparePassword(currUser.Password) {
		return 0, errors.New("wrong password")
	}
	
	return us.ID, nil
}

func (u *UserUsecase) GetRoles(user *model.User) ([]*model.Role, error) {
	currManager, err := u.getManagerByUserFromServer(user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "getManagerByUserFromServer()")
	}

	currCompany, err := u.getCompanyFromServer(currManager.CompanyId)
	if err != nil {
		return nil, errors.Wrap(err, "getCompanyFromServer()")
	}

	var roles []*model.Role

	// TODO: rewrite avatar in Role struct
	clientRole := &model.Role{
		Role:	"client",
		Label:	currCompany.CompanyName,
	}

	freelanceRole := &model.Role{
		Role:   "freelancer",
		Label:  user.FirstName + " " + user.SecondName,
	}

	roles = append(roles, clientRole)
	roles = append(roles, freelanceRole)

	return roles, nil
}