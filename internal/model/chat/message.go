package model

import (
	"time"
)

type Message struct {
	ID         	int64		`json:"id"`
	ChatID     	int64		`json:"chatId"`
	SenderID   	int64		`json:"senderId"`
	ReceiverID 	int64		`json:"receiverId"`
	Body       	string		`json:"body"`
	Date       	time.Time	`json:"date"`
	IsRead		bool		`json:"isRead"`
}

type Packet struct {
	Transaction	string		`json:"transaction"`
	Message		Message		`json:"message"`
	Chat		Chat		`json:"chat"`
	Client		bool		`json:"isClient,string"`
}