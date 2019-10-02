package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	ID 				int `json:"id"`
	FirstName 		string `json:"firstName"`
	SecondName 		string `json:"secondName"`
	UserName     	string `json:"username"`
	Email 			string `json:"email"`
	Password 		string `json:"-"`
	EncryptPassword string `json:"-"`
	Avatar bool `json:"-"`
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

type UserInput struct {
	Name 		string `json:"firstName"`
	Surname 	string `json:"secondName"`
	Email   	string `json:"email"`
	Password	string `json:"password"`
	UserType 	string `json:"type"`
}

func (u *UserInput) CheckEmail() error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if valid := re.MatchString(u.Email); valid == false {
		return errors.New("invalid email")
	}
	return nil
}
