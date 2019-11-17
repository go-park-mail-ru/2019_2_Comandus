package userUcase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

type UserUsecase struct {
	userRep			user.Repository
	managerRep		manager.Repository
	freelancerRep	freelancer.Repository
	companyRep		company.Repository
}

func NewUserUsecase(u user.Repository, m manager.Repository, f freelancer.Repository, c company.Repository) user.Usecase {
	return &UserUsecase{
		userRep:		u,
		managerRep:		m,
		freelancerRep:	f,
		companyRep:		c,
	}
}

func (usecase *UserUsecase) CreateUser(data *model.User) error {
	if err := data.Validate(); err != nil {
		return errors.Wrap(err, "CreateUser")
	}

	if err := data.BeforeCreate(); err != nil {
		return errors.Wrap(err, "CreateUser")
	}

	if err := usecase.userRep.Create(data); err != nil {
		return errors.Wrap(err, "CreateUser<-userRep.Create()")
	}

	c := &model.Company{}
	if err := usecase.companyRep.Create(c); err != nil {
		return errors.Wrap(err, "CreateUser<-companyRep.Create(): ")
	}

	m := &model.HireManager{
		AccountID:        data.ID,
		RegistrationDate: time.Now(),
		CompanyID:      	c.ID,
	}

	if err := usecase.managerRep.Create(m); err != nil {
		return errors.Wrap(err, "CreateUser<-managerRep.Create()")
	}

	f := &model.Freelancer{
		AccountId:         data.ID,
		RegistrationDate:  time.Now(),
	}

	if err := usecase.freelancerRep.Create(f); err != nil {
		return errors.Wrap(err, "CreateUser<-managerRep.Create()")
	}

	return nil
}

func (usecase * UserUsecase) EditUser(new *model.User, old * model.User) error {
	new.ID = old.ID

	if old.Email != new.Email {
		return errors.Wrap(errors.New("can't change email"), "EditUser")
	}

	new.UserType = old.UserType
	new.EncryptPassword = old.EncryptPassword

	if err := usecase.userRep.Edit(new); err != nil {
		return errors.Wrap(err, "userRep.Edit()")
	}
	return nil
}

func (usecase *UserUsecase) EditUserPassword(passwords *model.BodyPassword, user *model.User) error {
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

	if err := usecase.userRep.Edit(user); err != nil {
		return errors.Wrapf(err, "userRep.Edit")
	}
	return nil
}

func (usecase *UserUsecase) GetAvatar(user *model.User) ([]byte, error) {
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

func (usecase *UserUsecase) Find(id int64) (*model.User, error) {
	user, err := usecase.userRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}
	currFreelancer, err := usecase.freelancerRep.FindByUser(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	currManager, err := usecase.managerRep.FindByUser(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find()")
	}

	user.FreelancerId = currFreelancer.ID
	user.HireManagerId = currManager.ID
	user.CompanyId = currManager.CompanyID

	return user, nil
}

func (usecase *UserUsecase) SetUserType(user *model.User, userType string) error {
	if err := user.SetUserType(userType); err != nil {
		return errors.Wrap(err, "SetUserType()")
	}

	if err := usecase.userRep.Edit(user); err != nil {
		return errors.Wrap(err, "userRep.Edit()")
	}
	return nil
}

func (usecase *UserUsecase) VerifyUser(currUser *model.User) (int64, error) {
	u, err := usecase.userRep.FindByEmail(currUser.Email)
	if err != nil {
		return 0, errors.Wrapf(err, "userRep.FindByEmail()")
	}

	if !u.ComparePassword(currUser.Password) {
		return 0, errors.New("wrong password")
	}
	
	return u.ID, nil
}

func (usecase *UserUsecase) GetRoles(user *model.User) ([]*model.Role, error) {
	currManager, err := usecase.managerRep.FindByUser(user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "managerRep.FindByUser()")
	}

	currCompany, err := usecase.companyRep.Find(currManager.CompanyID)
	if err != nil {
		return nil, errors.Wrap(err, "companyRep.Find()")
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