package userRepository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"

)

type UserRepository struct {
	db	*sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(u *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (firstName, secondName, username, email, encryptPassword, userType) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING accountId",
		u.FirstName,
		u.SecondName,
		u.UserName,
		u.Email,
		u.EncryptPassword,
		u.UserType,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int64) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT accountId, firstName, secondName, username, email, '' as password, encryptPassword, avatar, userType FROM users WHERE accountId = $1",
		id,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.SecondName,
		&u.UserName,
		&u.Email,
		&u.Password,
		&u.EncryptPassword,
		&u.Avatar,
		&u.UserType,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.db.QueryRow(
		"SELECT accountId, firstName, secondName, username, email, '' as password, encryptPassword, avatar, userType FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.SecondName,
		&u.UserName,
		&u.Email,
		&u.Password,
		&u.EncryptPassword,
		&u.Avatar,
		&u.UserType,
	); err != nil {
		return nil, err
	}
	return u, nil
}

//TODO: validate user
func (r *UserRepository) Edit(u * model.User) error {
	return r.db.QueryRow("UPDATE users SET firstName = $1, secondName = $2, userName = $3, " +
		"encryptPassword = $4, avatar = $5, userType = $6 WHERE accountId = $7 RETURNING accountId",
		u.FirstName,
		u.SecondName,
		u.UserName,
		u.EncryptPassword,
		u.Avatar,
		u.UserType,
		u.ID,
	).Scan(&u.ID)
}