package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	userFreelancer = "freelancer"
	userCustomer   = "client"
	)

type User struct {
	ID 				int `json:"id"`
	FirstName 		string `json:"firstName"`
	SecondName 		string `json:"secondName"`
	UserName     	string `json:"username"`
	Email 			string `json:"email"`
	Password		string `json:"password"`
	EncryptPassword string `json:"-"`
	Avatar 			[]byte `json:"-"`
	UserType 		string `json:"type"`
}

func (u *User) BeforeCreate() error {
	if len(u.UserType) == 0 || u.UserType != userFreelancer && u.UserType != userCustomer {
		u.UserType = userFreelancer
	}

	if len(u.Password) > 0 {
		enc, err := EncryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptPassword = enc
	}
	u.Password = ""
	return nil
}

func (u *User) SetUserType(userType string) error {
	if userType == userFreelancer || userType == userCustomer {
		u.UserType = userType
		return nil
	}
	return errors.New("wrong user type")
}

func (u *User) IsManager() bool {
	return u.UserType == userCustomer
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

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}

