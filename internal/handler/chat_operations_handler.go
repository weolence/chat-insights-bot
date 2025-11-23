package handler

import (
	"context"
	"fmt"
	"log"
	"main/internal/controller"
	"main/internal/model"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/telebot.v3"
)

const (
	UnexpectedUserState  = "unexpected user state"
	UnexpectedFileFormat = "unexpected file format"
)

type ChatOperationsHandler struct {
	bot                *telebot.Bot
	chatController     *controller.ChatController
	llmController      *controller.LlmController
	userController     *controller.UserController
	btnRouteController *controller.BtnRouteController
}

func NewChatOperationsHandler(bot *telebot.Bot, chatController *controller.ChatController, llmController *controller.LlmController,
	userController *controller.UserController, btnRouteController *controller.BtnRouteController) *ChatOperationsHandler {
	return &ChatOperationsHandler{
		bot:                bot,
		chatController:     chatController,
		llmController:      llmController,
		userController:     userController,
		btnRouteController: btnRouteController,
	}
}

func (chatOperationsHandler *ChatOperationsHandler) SetupHandlers() {
	chatOperationsHandler.bot.Handle(telebot.OnText, chatOperationsHandler.handleText)
	chatOperationsHandler.bot.Handle(telebot.OnDocument, chatOperationsHandler.handleDocument)
	chatOperationsHandler.bot.Handle(telebot.OnCallback, chatOperationsHandler.handleRouteFunction)

	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnImportChat, chatOperationsHandler.handleImportChat)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnSelectChat, chatOperationsHandler.handleSelectChat)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnChatAnalyzation, chatOperationsHandler.handleChatAnalyzation)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnSummarizeChat, chatOperationsHandler.handleSummarizeChat)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnMeetingSearch, chatOperationsHandler.handleMeetingSearch)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnDescribePersonality, chatOperationsHandler.handleDescribePersonality)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnContextSearch, chatOperationsHandler.handleContextSearch)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnChatSettings, chatOperationsHandler.handleChatSettings)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnRenameChat, chatOperationsHandler.handleRenameChat)
	chatOperationsHandler.btnRouteController.SetRoute(controller.BtnRemoveChat, chatOperationsHandler.handleRemoveChat)
}

func (chatOperationsHandler *ChatOperationsHandler) handleRouteFunction(telebotCtx telebot.Context) error {
	data := telebotCtx.Callback().Data

	cleanData := strings.TrimSpace(data)

	routedFunc, ok := chatOperationsHandler.btnRouteController.GetRoute(cleanData)
	if !ok {
		fmt.Printf("found unexpected route: %s\n", data)
		return nil
	}

	telebotCtx.Respond()
	return routedFunc(telebotCtx)
}

func (chatOperationsHandler *ChatOperationsHandler) dropUserToRootMenu(telebotCtx telebot.Context, botError error) error {
	telegramUser := telebotCtx.Sender()

	menu := controller.CreateRootMenu()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)

	user.State = model.StateRootMenu
	user.SelectedChat = nil
	user.NewChatName = ""

	if err != nil {
		log.Println(err)
		return telebotCtx.Send("Возникла ошибка, попробуйте позже", menu)
	}

	if botError != nil {
		log.Println(botError)
		return telebotCtx.Send("Возникла ошибка, попробуйте позже", menu)
	}

	return telebotCtx.Send("Выберите действие:", menu)
}

// user pressed import chat
func (chatOperationsHandler *ChatOperationsHandler) handleImportChat(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateRootMenu {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	user.State = model.StateNameForChatAwaiting

	return telebotCtx.Send("Название чата не должно содержать ничего кроме цифр, латинских букв и пробелов. Введите название чата:")
}

// user sent text to bot
func (chatOperationsHandler *ChatOperationsHandler) handleText(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	switch user.State {
	case model.StateNameForChatAwaiting:
		name := telebotCtx.Text()
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\s]+$`, name)
		if err != nil || !ok {
			return telebotCtx.Send("Название не соответствует требованиям, попробуйте другое:")
		}

		user.State = model.StateFileOfChatAwaiting
		user.NewChatName = name

		return telebotCtx.Send("Прикрепите один файл чата в формате .html:")
	case model.StateDescriptionForContextSearchAwaiting:
		description := telebotCtx.Text()

		statusMsg, err := telebotCtx.Bot().Send(telebotCtx.Chat(), "Выполняется поиск контекста. Пожалуйста подождите...")
		if err != nil {
			return err
		}

		answer, err := chatOperationsHandler.llmController.ContextSearch(context.TODO(), *user.SelectedChat, description)
		if err != nil {
			return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
		}

		_, err = telebotCtx.Bot().Edit(statusMsg, fmt.Sprintf("Результаты поиска контекста:\n%s", answer))
		if err != nil {
			return err
		}

		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
	case model.StateNewNameForChatAwaiting:
		name := telebotCtx.Text()
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\s]+$`, name)
		if err != nil || !ok {
			return telebotCtx.Send("Название не соответствует требованиям, попробуйте другое:")
		}

		chatOperationsHandler.chatController.RenameChat(*user.SelectedChat, name)

		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
	default:
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}
}

// user sent file to bot
func (chatOperationsHandler *ChatOperationsHandler) handleDocument(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateFileOfChatAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	document := telebotCtx.Message().Document

	if strings.ToLower(filepath.Ext(document.FileName)) != ".html" {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedFileFormat))
	}

	reader, err := chatOperationsHandler.bot.File(&document.File)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	err = chatOperationsHandler.chatController.CreateChat(document.UniqueID, user.TelegramId, user.NewChatName, reader)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
}

// user pressed select chat
func (chatOperationsHandler *ChatOperationsHandler) handleSelectChat(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateRootMenu {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	importedChats, err := chatOperationsHandler.chatController.GetUserChats(*user)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}
	if len(importedChats) == 0 {
		telebotCtx.Send("Импортированных чатов не найдено")
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
	}

	menu, err := controller.CreateImportedChats(importedChats)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	user.State = model.StateChatSelectionAwaiting

	for _, chat := range importedChats {
		chatOperationsHandler.btnRouteController.SetRoute(fmt.Sprintf("chat_id_%s", chat.Id), chatOperationsHandler.handleChatButton)
	}

	return telebotCtx.Send("Выберите чат:", menu)
}

// user pressed one of chat buttons with chatId during chat selection
func (chatOperationsHandler *ChatOperationsHandler) handleChatButton(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatSelectionAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	chatId := strings.TrimPrefix(strings.TrimSpace(telebotCtx.Callback().Data), "chat_id_")

	importedChats, err := chatOperationsHandler.chatController.GetUserChats(*user)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	for _, chat := range importedChats {
		if chat.Id == chatId {
			user.SelectedChat = &model.Chat{
				Id:       chat.Id,
				UserId:   chat.UserId,
				Name:     chat.Name,
				Filepath: chat.Filepath,
			}
		}
		chatOperationsHandler.btnRouteController.DeleteRoute(fmt.Sprintf("chat_id_%s", chat.Id))
	}

	user.State = model.StateChatInteractionTypeAwaiting

	menu := controller.CreateAvailableChatInteractions()

	return telebotCtx.Send("Выберите тип взаимодействия с чатом:", menu)
}

// user pressed Analyze chat
func (chatOperationsHandler *ChatOperationsHandler) handleChatAnalyzation(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatInteractionTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	user.State = model.StateChatAnalyzationTypeAwaiting

	menu := controller.CreateAvailableAnalysisMethods()

	return telebotCtx.Send("Выберите метод анализа чата:", menu)
}

// user pressed summarize chat
func (chatOperationsHandler *ChatOperationsHandler) handleSummarizeChat(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatAnalyzationTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	statusMsg, err := telebotCtx.Bot().Send(telebotCtx.Chat(), "Выполняется анализ. Пожалуйста подождите...")
	if err != nil {
		return err
	}

	answer, err := chatOperationsHandler.llmController.SummarizeChat(context.TODO(), *user.SelectedChat)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	_, err = telebotCtx.Bot().Edit(statusMsg, fmt.Sprintf("Результаты подытоживания чата:\n%s", answer))
	if err != nil {
		return err
	}

	return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
}

// user pressed meeting search
func (chatOperationsHandler *ChatOperationsHandler) handleMeetingSearch(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatAnalyzationTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	statusMsg, err := telebotCtx.Bot().Send(telebotCtx.Chat(), "Выполняется поиск встреч. Пожалуйста подождите...")
	if err != nil {
		return err
	}

	answer, err := chatOperationsHandler.llmController.MeetingSearch(context.TODO(), *user.SelectedChat)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	_, err = telebotCtx.Bot().Edit(statusMsg, fmt.Sprintf("Результаты поиска встреч:\n%s", answer))
	if err != nil {
		return err
	}

	return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
}

// user pressed describe personality
func (chatOperationsHandler *ChatOperationsHandler) handleDescribePersonality(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatAnalyzationTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	statusMsg, err := telebotCtx.Bot().Send(telebotCtx.Chat(), "Выполняется анализ личности. Пожалуйста подождите...")
	if err != nil {
		return err
	}

	answer, err := chatOperationsHandler.llmController.DescribePersonality(context.TODO(), *user.SelectedChat)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	_, err = telebotCtx.Bot().Edit(statusMsg, fmt.Sprintf("Описание личности:\n%s", answer))
	if err != nil {
		return err
	}

	return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
}

// user pressed context search
func (chatOperationsHandler *ChatOperationsHandler) handleContextSearch(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatAnalyzationTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	user.State = model.StateDescriptionForContextSearchAwaiting

	return telebotCtx.Send("Опишите событие одним сообщением:")
}

// user pressed chat settings
func (chatOperationsHandler *ChatOperationsHandler) handleChatSettings(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatInteractionTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	user.State = model.StateChatSettingsTypeAwaiting

	menu := controller.CreateAvailableSettings()

	return telebotCtx.Send("Доступные настройки чата:", menu)
}

// user pressed rename chat
func (chatOperationsHandler *ChatOperationsHandler) handleRenameChat(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)
	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatSettingsTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	user.State = model.StateNewNameForChatAwaiting

	return telebotCtx.Send("Название чата не должно содержать ничего кроме цифр, латинских букв и пробелов. Введите новое название чата:")
}

// user pressed remove chat
func (chatOperationsHandler *ChatOperationsHandler) handleRemoveChat(telebotCtx telebot.Context) error {
	telegramUser := telebotCtx.Sender()

	user, err := chatOperationsHandler.userController.GetUser(telegramUser.ID)

	if err != nil || user == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	if user.State != model.StateChatSettingsTypeAwaiting {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	if user.SelectedChat == nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, fmt.Errorf(UnexpectedUserState))
	}

	err = chatOperationsHandler.chatController.RemoveChat(*user.SelectedChat)
	if err != nil {
		return chatOperationsHandler.dropUserToRootMenu(telebotCtx, err)
	}

	return chatOperationsHandler.dropUserToRootMenu(telebotCtx, nil)
}
