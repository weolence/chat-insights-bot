package repository

import (
	"database/sql"
	"main/internal/model"

	_ "modernc.org/sqlite"
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
			id TEXT PRIMARY KEY,
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

func (cr *ChatsRepository) CreateChat(chatId string, userId int64, name string, filepath string) error {
	_, err := cr.db.Exec(
		"INSERT INTO chats (id, user_id, name, filepath) VALUES (?, ?, ?, ?)",
		chatId, userId, name, filepath,
	)
	return err
}

func (cr *ChatsRepository) RemoveChat(chatId string) error {
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

func (cr *ChatsRepository) RenameChat(chatId string, newName string) error {
	_, err := cr.db.Exec(
		"UPDATE chats SET name = ? WHERE id = ?",
		newName, chatId,
	)
	return err
}
