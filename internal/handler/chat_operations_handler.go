package handler

import (
	"fmt"
	"main/internal/controller"
	"main/internal/model"
	"regexp"

	"gopkg.in/telebot.v3"
)

const (
	UserDataConversionErr = "user data not found"
	UnexpectedUserState   = "unexpected user state"
)

type ChatOperationsHandler struct {
	bot *telebot.Bot
	cc  *controller.ChatController
}

func NewChatOperationsHandler(bot *telebot.Bot, cc *controller.ChatController) *ChatOperationsHandler {
	return &ChatOperationsHandler{bot: bot, cc: cc}
}

func (coh *ChatOperationsHandler) SetupHandlers() {
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnRemoveChat}, coh.handleRemoveChat)
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnRenameChat}, coh.handleRenameChat)
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnSummarizeChat}, coh.handleSummarizeChat)
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnDescribePersonality}, coh.handleDescribePersonality)
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnMeetingSearch}, coh.handleMeetingSearch)
	coh.bot.Handle(&telebot.Callback{Data: controller.BtnContextSearch}, coh.handleContextSearch)
}

func dropUserToRootMenu(c telebot.Context, err error) error {
	telegramUser := c.Sender()
	model.UserStates.Store(telegramUser.ID, model.UserData{State: model.StateRootMenu})
	menu := controller.CreateRootMenu()
	if err != nil {
		return c.Send("Возникла ошибка, попробуйте позже", menu)
	}
	return c.Send("Выберите действие:", menu)
}

func getUserData(telegramUserId int64) (*model.UserData, bool) {
	value, ok := model.UserStates.Load(telegramUserId)
	if !ok {
		return nil, !ok
	}

	userData, ok := value.(model.UserData)
	if !ok {
		return nil, !ok
	}

	return &userData, ok
}

func (coh *ChatOperationsHandler) handleRemoveChat(c telebot.Context) error {
	telegramUser := c.Sender()
	userData, ok := getUserData(telegramUser.ID)
	if !ok {
		return dropUserToRootMenu(c, fmt.Errorf(UserDataConversionErr))
	}

	if userData.State != model.StateChatSettingsSelected {
		return dropUserToRootMenu(c, fmt.Errorf(UnexpectedUserState))
	}

	chat, err := coh.cc.GetChat(userData.SelectedChatId)
	if err != nil {
		return dropUserToRootMenu(c, err)
	}

	err = coh.cc.RemoveChat(*chat)
	if err != nil {
		return dropUserToRootMenu(c, err)
	}

	return dropUserToRootMenu(c, nil)
}

func (coh *ChatOperationsHandler) handleRenameChat(c telebot.Context) error {
	telegramUser := c.Sender()
	userData, ok := getUserData(telegramUser.ID)
	if !ok {
		return dropUserToRootMenu(c, fmt.Errorf(UserDataConversionErr))
	}

	if userData.State != model.StateChatSettingsSelected {
		return dropUserToRootMenu(c, fmt.Errorf(UnexpectedUserState))
	}

	model.UserStates.Store(telegramUser.ID, model.UserData{State: model.StateNewNameForChatAwaiting})

	return c.Send("Название чата не должно содержать ничего кроме цифр, латинских букв и пробелов. Введите новое название чата:")
}

func (coh *ChatOperationsHandler) handleText(c telebot.Context) error {
	telegramUser := c.Sender()
	userData, ok := getUserData(telegramUser.ID)
	if !ok {
		return dropUserToRootMenu(c, fmt.Errorf(UserDataConversionErr))
	}

	switch userData.State {
	case model.StateNewNameForChatAwaiting:
		newName := c.Text()
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\s]+$`, newName)
		if err != nil || !ok {
			return c.Send("Новое название не соответствует требованиям, попробуйте другое:")
		}

		err = coh.cc.RenameChat(userData.SelectedChatId, newName)
		if err != nil {
			return err
		}

		return c.Send("Название чата успешно изменено на " + newName)
	case model.StateNameForChatAwaiting:
		newName := c.Text()
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\s]+$`, newName)
		if err != nil || !ok {
			return c.Send("Название не соответствует требованиям, попробуйте другое:")
		}

	}
}

func (coh *ChatOperationsHandler) handleSummarizeChat(c telebot.Context) error {

}

func (coh *ChatOperationsHandler) handleDescribePersonality(c telebot.Context) error {

}

func (coh *ChatOperationsHandler) handleMeetingSearch(c telebot.Context) error {

}

func (coh *ChatOperationsHandler) handleContextSearch(c telebot.Context) error {

}
