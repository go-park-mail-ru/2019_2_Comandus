package chat_rep

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) chat.Repository {
	return &ChatRepository{db}
}

func (r *ChatRepository) Create(chat *model2.Chat) error {
	return r.db.QueryRow(
		"INSERT INTO chats (user_id, support_id, name, proposal_id) " +
			"VALUES ($1, $2, $3, $4) RETURNING id",
		chat.Freelancer,
		chat.Manager,
		chat.Name,
		chat.ProposalId,
	).Scan(&chat.ID)
}

func (r *ChatRepository) Find(id int64) (*model2.Chat, error) {
	c := &model2.Chat{}
	if err := r.db.QueryRow(
		"SELECT id, user_id, support_id, proposal_id, name FROM chats WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.Freelancer,
		&c.Manager,
		&c.ProposalId,
		&c.Name,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepository) FindByProposal(id int64) (*model2.Chat, error) {
	c := &model2.Chat{}
	if err := r.db.QueryRow(
		"SELECT id, user_id, support_id, proposal_id, name FROM chats WHERE proposal_id = $1;",
		id,
	).Scan(
		&c.ID,
		&c.Freelancer,
		&c.Manager,
		&c.ProposalId,
		&c.Name,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepository) FindCurrent(userId int64, clientId int64) (*model2.Chat, error) {
	c := &model2.Chat{}
	if err := r.db.QueryRow(
		"SELECT id, user_id, support_id, proposal_id, name FROM chats WHERE user_id = $1 AND support_id = $2;",
		userId, clientId,
	).Scan(
		&c.ID,
		&c.Freelancer,
		&c.Manager,
		&c.ProposalId,
		&c.Name,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepository) List(userId int64, clientId int64) ([]*model2.Chat, error) {
	var chats []*model2.Chat
	rows, err := r.db.Query(
		"SELECT id, user_id, support_id, proposal_id, name FROM chats " +
			"WHERE ($1 = 0 OR user_id = $1) AND " +
			"($2 = 0 OR support_id = $2) " +
			"ORDER BY id DESC;",
			userId, clientId,
			)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := &model2.Chat{}
		err := rows.Scan(&c.ID, &c.Freelancer, &c.Manager, &c.ProposalId, &c.Name)
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

