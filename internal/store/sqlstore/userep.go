package sqlstore

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"log"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	/*if err := u.Validate(); err != nil {
		return err
	}*/

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (firstName, secondName, username, email, encryptPassword) VALUES ($1, $2, $3, $4, $5) RETURNING accountId",
		u.FirstName,
		u.SecondName,
		u.UserName,
		u.Email,
		u.EncryptPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT accountId, email, encryptPassword FROM users WHERE accountId = $1",
		id,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.SecondName,
		&u.UserName,
		&u.Email,
		&u.EncryptPassword,
		&u.Avatar,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT accountId, email, encryptPassword FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptPassword,
	); err != nil {
		log.Println("hello find by email func", err)
		return nil, err
	}
	return u, nil
}
