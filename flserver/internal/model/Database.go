package model

import "sync"

type UsersDB struct {
	Users 			[]User
	Freelancers 	[]Freelancer
	HireManagers 	[]HireManager
	Companies	 	[]Company
	Notifications 	[]Notification
	InnerInfos		[]InnerInfo
	Mu    			*sync.Mutex
}

func NewUsersDB() *UsersDB {
	return &UsersDB{
		make([]User, 0),
		make([]Freelancer, 0),
		make([]HireManager, 0),
		make([]Company,0),
		make([]Notification,0),
		make([]InnerInfo,0),
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

func (db *UsersDB) GetCompanyByID (id int) *Company {
	for i := 0; i < len(db.Companies); i++ {
		if id == db.Companies[i].ID {
			return &db.Companies[i]
		}
	}
	return nil
}

func (db *UsersDB) GetFreelancerByID (id int) *Freelancer {
	for i := 0; i < len(db.Freelancers); i++ {
		if id == db.Freelancers[i].ID {
			return &db.Freelancers[i]
		}
	}
	return nil
}

//func (db *UsersDB) GetProfileByID (id int64) *Profile {
//	for i := 0; i < len(db.Profiles); i++ {
//		if id == db.Profiles[i].ID {
//			return &db.Profiles[i]
//		}
//	}
//	return nil
//}



