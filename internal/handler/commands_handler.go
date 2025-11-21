package handler

import (
	"fmt"
	"main/internal/controller"

	"gopkg.in/telebot.v3"
)

type CommandsHandler struct {
	bot            *telebot.Bot
	userController *controller.UserController
}

func NewCommandsHandler(bot *telebot.Bot, userController *controller.UserController) *CommandsHandler {
	return &CommandsHandler{bot: bot, userController: userController}
}

func (commandsHandler *CommandsHandler) SetupHandlers() {
	commandsHandler.bot.Handle("/start", commandsHandler.handleStartMessage)
}

func (commandsHandler *CommandsHandler) handleStartMessage(c telebot.Context) error {
	telegramUser := c.Sender()

	user, err := commandsHandler.userController.GetUser(telegramUser.ID)
	if err != nil {

		return c.Send("Возникла ошибка во время регистрации, попробуйте позже")
	}

	menu := controller.CreateRootMenu()

	if user != nil {
		return c.Send("Вы уже начали использование бота. Выберите действие:", menu)
	}

	_, err = commandsHandler.userController.CreateUser(telegramUser.ID)
	if err != nil {
		fmt.Println(err)
		return c.Send("Возникла ошибка во время регистрации, попробуйте позже")
	}

	return c.Send(
		"Привет! Этот бот может анализировать твои чаты. "+
			"Всё что нужно — экспортировать чат в формате HTML и загрузить его сюда.",
		menu,
	)
}
