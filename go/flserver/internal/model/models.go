package model

import "sync"

type User struct {
	ID       int64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
	EncryptPassword string `json:"-"`
}

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UsersDB struct {
	Users []User
	Mu    *sync.Mutex
}

func NewUsersDB() *UsersDB {
	return &UsersDB{make([]User, 0), &sync.Mutex{} }
}