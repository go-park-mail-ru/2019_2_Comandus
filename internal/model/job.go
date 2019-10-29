package model

type Job struct {
	ID                int     `json:"id" valid:"int, optional"`
	HireManagerId     int     `json:"hireManagerId,string" valid:"int, optional"`
	Title             string  `json:"title" valid:"utfletternum, required"`
	Description       string  `json:"description"valid:"- , optional"`
	Files             string  `json:"files" valid:"-"`
	SpecialityId      int     `json:"specialityId,string" valid:"int, optional"`
	ExperienceLevelId int     `json:"experienceLevelId,string" valid:"in(1|2|3)"`
	PaymentAmount      float64 `json:"paymentAmount,string" valid:"float"`
	Country           string  `json:"country" valid:"utfletternum, optional"`
	City              string  `json:"city" valid:"utfletternum, optional"`
	JobTypeId         int     `json:"jobTypeId,string" valid:"int, optional"`
}

//curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "USANews", "description" : "bbbbbbb",
//"country" : "USA"}'
// curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "RussianNews", "description" : "aaaaaaa",
// "country" : "russia"}'
// curl -XPOST -v -b cookie.txt http://127.0.0.1:8080/jobs --data '{"title" : "ArmeniaNews", "description" : "armyDeth",
// "country" : "armenia"}'