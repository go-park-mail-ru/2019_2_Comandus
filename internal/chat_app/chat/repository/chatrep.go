package chat_rep

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model/chat"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) chat.Repository {
	return &ChatRepository{db}
}

func (r *ChatRepository) Create(chat *model.Chat) error {
	return r.db.QueryRow(
		"INSERT INTO chats (user_id, support_id, name) " +
			"VALUES ($1, $2, $3) RETURNING id",
		chat.UserID,
		chat.SupportID,
		chat.Name,
	).Scan(&chat.ID)
}

func (r *ChatRepository) Find(id int64) (*model.Chat, error) {
	c := &model.Chat{}
	if err := r.db.QueryRow(
		"SELECT id, user_id, support_id, name FROM chats WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.SupportID,
		&c.Name,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepository) FindCurrent(userId int64, clientId int64) (*model.Chat, error) {
	c := &model.Chat{}
	if err := r.db.QueryRow(
		"SELECT id, user_id, support_id, name FROM chats WHERE user_id = $1 AND support_id = $2;",
		userId, clientId,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.SupportID,
		&c.Name,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepository) List(userId int64, clientId int64) ([]*model.Chat, error) {
	var chats []*model.Chat
	rows, err := r.db.Query(
		"SELECT id, user_id, support_id, name FROM chats " +
			"WHERE ($1 = 0 OR user_id = $1) AND " +
			"($2 = 0 OR support_id = $2) " +
			"ORDER BY id DESC;",
			userId, clientId,
			)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := &model.Chat{}
		err := rows.Scan(&c.ID, &c.UserID, &c.SupportID, &c.Name)
		if err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return chats, nil
}

