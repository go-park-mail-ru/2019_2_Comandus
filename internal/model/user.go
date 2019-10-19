package model

import (
	"github.com/go-ozzo/ozzo-validation/is"
	//"errors"
	"golang.org/x/crypto/bcrypt"
	//"regexp"
	validation "github.com/go-ozzo/ozzo-validation"
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

