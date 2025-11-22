package controller

import (
	"log"
	"main/internal/model"
	"main/internal/repository"

	"github.com/cornelk/hashmap"
)

type UserController struct {
	usersRepo *repository.UsersRepository
	currUsers *hashmap.Map[int64, *model.User]
}

func NewUserController() (*UserController, error) {
	usersRepo, err := repository.NewUsersRepository()
	if err != nil {
		return nil, err
	}

	return &UserController{
		usersRepo: usersRepo,
		currUsers: hashmap.New[int64, *model.User](),
	}, nil
}

// creates new user in database and in table of current(online) users
func (uc *UserController) CreateUser(telegramId int64) (*model.User, error) {
	registered, err := uc.usersRepo.IsUserRegistered(telegramId)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		TelegramId: telegramId,
		State:      model.StateRootMenu,
	}

	if !registered {
		err = uc.usersRepo.CreateUser(telegramId)
		if err != nil {
			return nil, err
		}
	}

	uc.currUsers.Set(telegramId, user)

	return user, nil
}

/*
if user not registered in database nil will be returned,
in case user registered in database, but not in table
of current users, User struct will be fetched in it and returned
*/
func (uc *UserController) GetUser(telegramId int64) (*model.User, error) {
	user, ok := uc.currUsers.Get(telegramId)
	log.Println(user, ok)
	if ok {
		return user, nil
	}

	registered, err := uc.usersRepo.IsUserRegistered(telegramId)
	if err != nil {
		return nil, err
	}
	if !registered {
		return nil, nil
	}

	user = &model.User{
		TelegramId: telegramId,
		State:      model.StateRootMenu,
	}

	uc.currUsers.Set(telegramId, user)

	return user, nil
}
