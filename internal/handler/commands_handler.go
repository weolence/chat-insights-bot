package handler

import (
	"main/internal/controller"
	"main/internal/model"

	"gopkg.in/telebot.v3"
)

type CommandsHandler struct {
	bot *telebot.Bot
	uc  *controller.UserController
}

func NewCommandsHandler(bot *telebot.Bot, uc *controller.UserController) *CommandsHandler {
	return &CommandsHandler{bot: bot, uc: uc}
}

func (ch *CommandsHandler) SetupHandlers() {
	ch.bot.Handle("/start", ch.handleStartMessage)
}

func (ch *CommandsHandler) handleStartMessage(c telebot.Context) error {
	telegramUser := c.Sender()

	user, err := ch.uc.GetUser(telegramUser.ID)
	if err != nil {
		return c.Send("Возникла ошибка во время регистрации, попробуйте позже")
	}

	menu := controller.CreateRootMenu()
	if user != nil {
		model.UserStates.Store(telegramUser.ID, model.StateRootMenu)
		return c.Send("Вы уже начали использование бота. Выберите действие:", menu)
	}

	model.UserStates.Store(telegramUser.ID, model.StateRootMenu)
	return c.Send(
		"Привет! Этот бот может анализировать твои чаты. "+
			"Всё что нужно — экспортировать чат в формате HTML и загрузить его сюда.",
		menu,
	)
}
