package model

type Customer struct {
	ID       int `json:"id"`
	Company string `json:"company"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}
