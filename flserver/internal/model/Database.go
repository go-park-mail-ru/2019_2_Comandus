package model

import "sync"

type UsersDB struct {
	Users []User
	Freelancers []Freelancer
	Customers []Customer
	Profiles []Profile
	Mu    *sync.Mutex
}

func NewUsersDB() *UsersDB {
	return &UsersDB{
		make([]User, 0),
		make([]Freelancer, 0),
		make([]Customer, 0),
		make([]Profile,0),
		&sync.Mutex{},
	}
}

func (db *UsersDB) GetUserByID (id int) *User {
	for i := 0; i < len(db.Users); i++ {
		if id == db.Users[i].ID {
			return &db.Users[i]
		}
	}
	return nil
}

func (db *UsersDB) GetProfileByID (id int) *Profile {
	for i := 0; i < len(db.Profiles); i++ {
		if id == db.Profiles[i].ID {
			return &db.Profiles[i]
		}
	}
	return nil
}



