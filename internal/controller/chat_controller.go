package controller

import (
	"fmt"
	"io"
	"main/internal/model"
	"main/internal/repository"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/telebot.v3"
)

const (
	FileFormat      = "html"
	chatsStorageDir = "chats_storage"
)

type ChatController struct {
	chatsRepo *repository.ChatsRepository
}

func NewChatController() (*ChatController, error) {
	chatsRepo, err := repository.NewChatsRepository()
	if err != nil {
		return nil, err
	}
	return &ChatController{chatsRepo: chatsRepo}, nil
}

func (cc *ChatController) CreateChat(chatId string, userId int64, name string, file telebot.File) error {
	resp, err := http.Get(file.FileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	path := fmt.Sprintf("./%s/%d/%s.%s", chatsStorageDir, userId, file.UniqueID, FileFormat)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	err = cc.chatsRepo.CreateChat(chatId, userId, name, path)
	return err
}

func (cc *ChatController) RemoveChat(chat model.Chat) error {
	err := os.Remove(chat.Filepath)
	if err != nil {
		return err
	}

	return cc.chatsRepo.RemoveChat(chat.Id)
}

func (cc *ChatController) RenameChat(chat model.Chat, newName string) error {
	return cc.chatsRepo.RenameChat(chat.Id, newName)
}

func (cc *ChatController) GetUserChats(user model.User) ([]model.Chat, error) {
	return cc.chatsRepo.GetUserChats(user.TelegramId)
}
