package controller

import (
	"main/internal/model"
	"main/internal/repository"
)

type UserController struct {
	usersRepo *repository.UsersRepository
}

func NewUserController() (*UserController, error) {
	usersRepo, err := repository.NewUsersRepository()
	if err != nil {
		return nil, err
	}
	return &UserController{usersRepo: usersRepo}, nil
}

func (uc *UserController) CreateUser(telegramId int64) error {
	return uc.usersRepo.AddUser(telegramId)
}

func (uc *UserController) GetUser(telegramId int64) (*model.User, error) {
	return uc.usersRepo.GetUser(telegramId)
}
