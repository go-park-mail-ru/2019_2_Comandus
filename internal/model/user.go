package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	UserFreelancer = "freelancer"
	UserCustomer   = "client"
)

type User struct {
	ID 				int64 `json:"-" valid:"int, optional"`
	FirstName 		string `json:"firstName" valid:"utfletter, required"`
	SecondName 		string `json:"secondName" valid:"utfletter"`
	UserName     	string `json:"username" valid:"alphanum"`
	Email 			string `json:"email" valid:"email"`
	Password		string `json:"password" valid:"length(6|100)"`
	EncryptPassword string `json:"-" valid:"-"`
	Avatar 			[]byte `json:"-" valid:"-"`
	UserType 		string `json:"type" valid:"in(client|freelancer)"`
	RegistrationDate  	time.Time `json:"registrationDate" valid:"-"`
	FreelancerId    int64  `json:"freelancerId" valid:"int, optional"`
	HireManagerId   int64  `json:"hireManagerId" valid:"int, optional"`
	CompanyId       int64  `json:"companyId" valid:"int, optional"`
}

func (u *User) BeforeCreate() error {
	if len(u.UserType) == 0 || u.UserType != UserFreelancer && u.UserType != UserCustomer {
		u.UserType = UserFreelancer
	}

	if len(u.Password) > 0 {
		enc, err := EncryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptPassword = enc
	}

	u.Password = ""
	u.RegistrationDate = time.Now()

	return nil
}

func (u *User) SetUserType(userType string) error {
	if userType == UserFreelancer || userType == UserCustomer {
		u.UserType = userType
		return nil
	}
	return errors.New("wrong user type")
}

func (u *User) IsManager() bool {
	return u.UserType == UserCustomer
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptPassword), []byte(password)) == nil
}

func EncryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 100)),
	)
}

func (u *User) Sanitize(sanitizer *bluemonday.Policy) {
	u.FirstName = sanitizer.Sanitize(u.FirstName)
	u.SecondName = sanitizer.Sanitize(u.SecondName)
	u.UserName = sanitizer.Sanitize(u.UserName)
}
