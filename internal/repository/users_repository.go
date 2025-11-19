package repository

import (
	"database/sql"
	"main/internal/model"
)

type UsersRepository struct {
	db *sql.DB
}

const (
	usersDbPath = "./database/users_repository.db"
)

func NewUsersRepository() (*UsersRepository, error) {
	db, err := sql.Open("sqlite", usersDbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			telegram_id BIGINT PRIMARY KEY
		);
	`)
	if err != nil {
		return nil, err
	}

	return &UsersRepository{db: db}, nil
}

func (ur *UsersRepository) AddUser(telegramId int64) error {
	_, err := ur.db.Exec(
		"INSERT INTO users (telegram_id) VALUES (?)",
		telegramId,
	)
	return err
}

func (ur *UsersRepository) GetUser(telegramId int64) (*model.User, error) {
	row := ur.db.QueryRow(
		"SELECT telegram_id FROM users WHERE telegram_id = ?",
		telegramId,
	)

	var u model.User
	err := row.Scan(&u.TelegramId)

	switch err {
	case nil:
		return &u, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
