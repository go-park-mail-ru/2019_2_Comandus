package model

type Customer struct {
	ID       int64 `json:"id"`
	Company string `json:"company"`
	Address string `json:"address"`
	Phone string `json:"phone"`
}
