package insightsbot

import (
	"fmt"
	"main/internal/controller"
	"main/internal/handler"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

type InsightsBot struct {
	ch  *handler.CommandsHandler
	coh *handler.ChatOperationsHandler
	bot *telebot.Bot
}

func NewInsightsBot() (*InsightsBot, error) {
	botKey, ok := os.LookupEnv("CHATINSIGHTS_BOT_KEY")
	if !ok {
		return nil, fmt.Errorf("environment variable OPENAI_API_KEY required")
	}

	pref := telebot.Settings{
		Token:  botKey,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	userController, err := controller.NewUserController()
	if err != nil {
		return nil, err
	}

	commandsHandler := handler.NewCommandsHandler(bot, userController)

	chatController, err := controller.NewChatController()
	if err != nil {
		return nil, err
	}

	llmController, err := controller.NewLlmController()
	if err != nil {
		return nil, err
	}

	btnRouteController := controller.NewBtnRouteController()

	chatOperationsHandler := handler.NewChatOperationsHandler(bot, chatController, llmController, userController, btnRouteController)

	return &InsightsBot{commandsHandler, chatOperationsHandler, bot}, nil
}

func (ib *InsightsBot) Run() {
	ib.setupHandlers()
	ib.bot.Start()
}

func (ib *InsightsBot) setupHandlers() {
	ib.ch.SetupHandlers()
	ib.coh.SetupHandlers()
}
