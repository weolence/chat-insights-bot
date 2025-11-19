package repository

import (
	"database/sql"
	"main/internal/model"
)

type ChatsRepository struct {
	db *sql.DB
}

const (
	chatsDbPath = "./database/chats_repository.db"
)

func NewChatsRepository() (*ChatsRepository, error) {
	db, err := sql.Open("sqlite", chatsDbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			id BIGINT PRIMARY KEY AUTOINCREMENT,
			user_id BIGINT NOT NULL,
			name TEXT NOT NULL,
			filepath TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	return &ChatsRepository{db: db}, nil
}

func (cr *ChatsRepository) AddChat(userId int64, name string, filepath string) error {
	_, err := cr.db.Exec(
		"INSERT INTO chats (user_id, name, filepath) VALUES (?, ?, ?)",
		userId, name, filepath,
	)
	return err
}

func (cr *ChatsRepository) RemoveChat(chatId int64) error {
	_, err := cr.db.Exec(
		"DELETE FROM chats WHERE id = ?",
		chatId,
	)
	return err
}

func (cr *ChatsRepository) GetUserChats(userId int64) ([]model.Chat, error) {
	rows, err := cr.db.Query(
		"SELECT id, user_id, name, filepath FROM chats WHERE user_id = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat

	for rows.Next() {
		var chat model.Chat
		err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name, &chat.Filepath)
		switch err {
		case nil:
		case sql.ErrNoRows:
			return chats, nil
		default:
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (cr *ChatsRepository) GetChat(chatId int64) (*model.Chat, error) {
	row := cr.db.QueryRow(
		"SELECT id, user_id, name, filepath FROM chats WHERE id = ?",
		chatId,
	)

	var chat model.Chat
	err := row.Scan(&chat.Id, &chat.UserId, &chat.Name, &chat.Filepath)

	switch err {
	case nil:
		return &chat, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (cr *ChatsRepository) RenameChat(chatId int64, newName string) error {
	_, err := cr.db.Exec(
		"UPDATE chats SET name = ? WHERE id = ?",
		newName, chatId,
	)
	return err
}
