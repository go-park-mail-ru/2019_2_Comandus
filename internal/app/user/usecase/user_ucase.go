package usecase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

type Usecase struct {
	userRep			user.Repository
	managerRep		manager.Repository
	freelancerRep	freelancer.Repository
}

func NewUserUsecase(u user.Repository, m manager.Repository, f freelancer.Repository) user.Usecase {
	return &Usecase{
		userRep:		u,
		managerRep:		m,
		freelancerRep:	f,
	}
}

func (usecase *Usecase) CreateUser(data *model.User) error {
	if err := data.Validate(); err != nil {
		return errors.Wrap(err, "CreateUser: ")
	}

	if err := data.BeforeCreate(); err != nil {
		return errors.Wrap(err, "CreateUser: ")
	}

	if err := usecase.userRep.Create(data); err != nil {
		return errors.Wrap(err, "CreateUser<-userRep.Create(): ")
	}

	m := &model.HireManager{
		AccountID:        data.ID,
		RegistrationDate: time.Now(),
		CompanyID:        0,		//TODO: set default company
	}

	if err := usecase.managerRep.Create(m); err != nil {
		return errors.Wrap(err, "CreateUser<-managerRep.Create(): ")
	}

	f := &model.Freelancer{
		AccountId:         data.ID,
		RegistrationDate:  time.Now(),
	}

	if err := usecase.freelancerRep.Create(f); err != nil {
		return errors.Wrap(err, "CreateUser<-managerRep.Create(): ")
	}

	return nil
}

func (usecase * Usecase) EditUser(new *model.User, old * model.User) error {
	if old.ID != new.ID {
		return errors.Wrap(errors.New("wrong user ID"), "EditUser: ")
	}

	if old.Email != new.Email {
		return errors.Wrap(errors.New("can't change email"), "EditUser: ")
	}

	if !old.ComparePassword(new.Password) {
		return errors.Wrap(errors.New("can't change password without validation"),
			"ComparePassword: ")
	}

	if err := usecase.userRep.Edit(new); err != nil {
		return errors.Wrap(err, "userRep.Edit(): ")
	}
	return nil
}

func (usecase *Usecase) EditUserPassword(passwords *model.BodyPassword, user *model.User) error {
	if passwords.NewPassword != passwords.NewPasswordConfirmation {
		return errors.New("new passwords are different")
	}

	if !user.ComparePassword(passwords.Password) {
		err := errors.New("wrong old password")
		return errors.Wrapf(err, "model.user.ComparePassword: ")
	}

	newEncryptPassword, err := model.EncryptString(passwords.NewPassword)
	if err != nil {
		return errors.Wrap(err,"model.EncryptString: ")
	}
	user.EncryptPassword = newEncryptPassword

	if err := usecase.userRep.Edit(user); err != nil {
		return errors.Wrapf(err, "userRep.Edit: ")
	}
	return nil
}

func (usecase *Usecase) GetAvatar(user *model.User) ([]byte, error) {
	if user.Avatar != nil {
		return user.Avatar, nil
	}

	var openFile *os.File
	filename := "internal/store/avatars/default.png"
	openFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Open: ")
	}

	defer func() {
		if err := openFile.Close(); err != nil {
			// TODO: write in correct logger
			log.Println(errors.Wrap(err, "GetAvatar<-Close(): "))
		}
	}()

	avatar := make([]byte, 0)
	if _, err := openFile.Read(avatar); err != nil {
		return nil, errors.Wrap(err, "Read(): ")
	}

	return avatar, nil
}

func (usecase *Usecase) Find(id int64) (*model.User, error) {
	user, err := usecase.userRep.Find(id)
	if err != nil {
		return nil, errors.Wrap(err, "userRep.Find(): ")
	}
	return user, nil
}

func (usecase *Usecase) SetUserType(user *model.User, userType string) error {
	if err := user.SetUserType(userType); err != nil {
		return errors.Wrap(err, "SetUserType(): ")
	}

	if err := usecase.userRep.Edit(user); err != nil {
		return errors.Wrap(err, "userRep.Edit(): ")
	}
	return nil
}

func (usecase *Usecase) VerifyUser(currUser *model.User) (int64, error) {
	u, err := usecase.userRep.FindByEmail(currUser.Email)
	if err != nil {
		return 0, errors.Wrapf(err, "userRep.FindByEmail(): ")
	}

	if !u.ComparePassword(currUser.Password) {
		return 0, errors.New("wrong password")
	}
	
	return u.ID, nil
}