package userUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
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
	_, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "os.Getwd()")
	}

	filename := "../internal/store/avatars/default.png"
	openFile, err = os.Open(filename)
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

func (u *UserUsecase) Find(id int64) (*model.User, error) {
	user, err := u.userRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	currFreelancer, err := clients.GetFreelancerByUserFromServer(id)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetFreelancerByUserFromServer()")
	}

	currManager, err := clients.GetManagerByUserFromServer(id)
	if err != nil {
		return nil, errors.Wrap(err, "clients.GetManagerByUserFromServer()")
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
	currManager, err := clients.GetManagerByUserFromServer(user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "clients.GetManagerByUserFromServer()")
	}

	currCompany, err := clients.GetCompanyFromServer(currManager.CompanyId)
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