package repository

import (
	"database/sql"
	"main/internal/model"

	_ "modernc.org/sqlite"
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

func (ur *UsersRepository) CreateUser(telegramId int64) error {
	_, err := ur.db.Exec(
		"INSERT INTO users (telegram_id) VALUES (?)",
		telegramId,
	)
	return err
}

func (ur *UsersRepository) IsUserRegistered(telegramId int64) (bool, error) {
	row := ur.db.QueryRow(
		"SELECT telegram_id FROM users WHERE telegram_id = ?",
		telegramId,
	)

	var u model.User
	err := row.Scan(&u.TelegramId)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
