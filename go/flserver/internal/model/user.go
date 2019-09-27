package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"sync"
)

type User struct {
	ID int64 `json:"id"`
	FreelancerID int64 `json:"freelancer"`
	CustomerID int64 `json:"customer"`
	Name     string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password, omitempty"`
	EncryptPassword string `json:"-"`
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptPassword = enc
	}
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type UserInput struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email     string `json:"email"`
	Password string `json:"password"`
}

func (u *UserInput) CheckEmail() error{
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if valid := re.MatchString(u.Email); valid == false {
		return errors.New("invalid email")
	}
	return nil
}

type UsersDB struct {
	Users []User
	Freelancers []Freelancer
	Customers []Customer
	Mu    *sync.Mutex
}

func NewUsersDB() *UsersDB {
	return &UsersDB{
		make([]User, 0),
		make([]Freelancer, 0),
		make([]Customer, 0),
		&sync.Mutex{},
	}
}

