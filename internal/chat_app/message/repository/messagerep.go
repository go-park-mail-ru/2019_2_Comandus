package mes_rep

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/message"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"log"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) message.Repository {
	return &MessageRepository{db}
}

func (r *MessageRepository) Create(message *model.Message) error {
	return r.db.QueryRow(
		"INSERT INTO messages (chat_id, sender_id, receiver_id, message, date, is_read) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		message.ChatID,
		message.SenderID,
		message.ReceiverID,
		message.Body,
		message.Date,
		message.IsRead,
	).Scan(&message.ID)
}

func (r *MessageRepository) Edit(message *model.Message) error {
	return r.db.QueryRow("UPDATE messages SET message = $1 WHERE id = $2 RETURNING id",
		message.Body,
		message.ID,
	).Scan(&message.ID)
}

func (r *MessageRepository) UpdateStatus(chatId int64, userId int64) error {
	log.Println("update status chat id: ", chatId, ", userId: ", userId)
	rows, err := r.db.Query("UPDATE messages SET is_read = $1 WHERE chat_id = $2 AND " +
		"receiver_id = $3 AND " +
		"NOT is_read RETURNING id, chat_id, sender_id, receiver_id, message, date, is_read;",
		true,
		chatId,
		userId,
	)
	log.Println(rows)
	return err
}

func (r *MessageRepository) List(chatId int64) ([]*model.Message, error) {
	var messages []*model.Message
	rows, err := r.db.Query(
		"SELECT id, chat_id, sender_id, receiver_id, message, date, is_read FROM messages " +
			"WHERE chat_id = $1 " +
			"ORDER BY date DESC", chatId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		m := &model.Message{}
		err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.ReceiverID, &m.Body, &m.Date, &m.IsRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}