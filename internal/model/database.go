package model

import (
	"sync"
)

type UsersDB struct {
	Users 			map[int] User
	Freelancers 	map[int] Freelancer
	HireManagers 	map[int] HireManager
	Companies	 	map[int] Company
	Notifications 	map[int] Notification
	InnerInfos		map[int] InnerInfo
	Jobs			map[int] Job
	Mu    			*sync.Mutex
	//ImageStore 	map[int] image.Image
	ImageStore 		map[int] []byte
}

func NewUsersDB() *UsersDB {
	return &UsersDB{
		make(map[int]User, 0),
		make(map[int]Freelancer, 0),
		make(map[int]HireManager, 0),
		make(map[int]Company,0),
		make(map[int]Notification,0),
		make(map[int]InnerInfo,0),
		make(map[int]Job,0),
		&sync.Mutex{},
		make(map[int][]byte),
	}
}

/*func (db *UsersDB) GetUserByID (id int) *User {
	for i := 0; i < len(db.Users); i++ {
		if id == db.Users[i].ID {
			return &db.Users[i]
		}
	}
	return nil
}*/

/*func (db *UsersDB) GetCompanyByID (id int) *Company {
	for i := 0; i < len(db.Companies); i++ {
		if id == db.Companies[i].ID {
			return &db.Companies[i]
		}
	}
	return nil
}

func (db *UsersDB) GetFreelancerByID (id int) (*Freelancer, error) {
	for i := 0; i < len(db.Freelancers); i++ {
		if id == db.Freelancers[i].ID {
			return &db.Freelancers[i] , nil
		}
	}
	return nil , fmt.Errorf("no profile with this id")
}

func (db *UsersDB) GetFreelancerByUserID (id int) *Freelancer {
	for i := 0; i < len(db.Freelancers); i++ {
		if id == db.Freelancers[i].AccountId {
			return &db.Freelancers[i]
		}
	}
	return nil
}




func (db *UsersDB) GetInnerInfoByUserID (id int) *InnerInfo {
	for i := 0; i < len(db.InnerInfos); i++ {
		if id == db.InnerInfos[i].UserID {
			return &db.InnerInfos[i]
		}
	}
	return nil
}

func (db *UsersDB) GetNotificationsByUserID (id int) *Notification {
	for i := 0; i < len(db.Notifications); i++ {
		if id == db.Notifications[i].UserID {
			return &db.Notifications[i]
		}
	}
	return nil
}

func (db *UsersDB) GetHireManagerByID (id int) *HireManager {
	for i := 0; i < len(db.HireManagers); i++ {
		if id == db.HireManagers[i].ID {
			return &db.HireManagers[i]
		}
	}
	return nil
}*/



