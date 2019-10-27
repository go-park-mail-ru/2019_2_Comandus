package model

type Job struct {
	ID                int     `json:"id"`
	HireManagerId     int     `json:"hireManagerId,string"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Files             string  `json:"files"`
	SpecialityId      int     `json:"specialityId,string"`
	ExperienceLevelId int     `json:"experienceLevelId,string"`
	PaymentAmount      float64 `json:"paymentAmount,string"`
	Country           string  `json:"country"`
	City              string  `json:"city"`
	JobTypeId         int     `json:"jobTypeId,string"`
}

//curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "USANews", "description" : "bbbbbbb",
//"country" : "USA"}'
// curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "RussianNews", "description" : "aaaaaaa",
// "country" : "russia"}'
// curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "ArmeniaNews", "description" : "armyDeth",
// "country" : "armenia"}'