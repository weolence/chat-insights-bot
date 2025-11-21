package controller

import (
	"fmt"
	"main/internal/model"

	"gopkg.in/telebot.v3"
)

const (
	BtnImportChat          = "import_chat_btn"
	BtnSelectChat          = "select_chat_btn"
	BtnChatAnalyzation     = "chat_analyzation_btn"
	BtnChatSettings        = "chat_settings_btn"
	BtnStopAnalyzing       = "stop_analyzing_btn"
	BtnSummarizeChat       = "summarize_chat_btn"
	BtnDescribePersonality = "describe_personality_btn"
	BtnMeetingSearch       = "meeting_search_btn"
	BtnContextSearch       = "context_search_btn"
	BtnRenameChat          = "rename_chat_btn"
	BtnRemoveChat          = "remove_chat_btn"
)

func CreateRootMenu() (menu *telebot.ReplyMarkup) {
	menu = &telebot.ReplyMarkup{}

	btnImportChat := menu.Data("Import chat", BtnImportChat, "")
	btnSelectChat := menu.Data("Select chat", BtnSelectChat, "")

	menu.Inline(
		menu.Row(btnImportChat),
		menu.Row(btnSelectChat),
	)

	return
}

func CreateAvailableChatInteractions() (menu *telebot.ReplyMarkup) {
	menu = &telebot.ReplyMarkup{}

	btnAnalyzeChat := menu.Data("Analyze chat", BtnChatAnalyzation)
	btnChatSettings := menu.Data("Chat settings", BtnChatSettings)

	menu.Inline(
		menu.Row(btnAnalyzeChat),
		menu.Row(btnChatSettings),
	)

	return
}

func CreateAvailableAnalysisMethods() (menu *telebot.ReplyMarkup) {
	menu = &telebot.ReplyMarkup{}

	btnStopAnalyzing := menu.Data("Stop analyzing", BtnStopAnalyzing)
	btnSummarizeChat := menu.Data("Summarize", BtnSummarizeChat)
	btnDescribePersonality := menu.Data("Describe personality", BtnDescribePersonality)
	btnMeetingSearch := menu.Data("Meeting search", BtnMeetingSearch)
	btnContextSearch := menu.Data("Context search", BtnContextSearch)

	menu.Inline(
		menu.Row(btnStopAnalyzing),
		menu.Row(btnSummarizeChat),
		menu.Row(btnDescribePersonality),
		menu.Row(btnMeetingSearch),
		menu.Row(btnContextSearch),
	)

	return
}

func CreateAvailableSettings() (menu *telebot.ReplyMarkup) {
	menu = &telebot.ReplyMarkup{}

	btnChangeName := menu.Data("Rename chat", BtnRenameChat)
	btnRemoveChat := menu.Data("Remove chat", BtnRemoveChat)

	menu.Inline(
		menu.Row(btnChangeName),
		menu.Row(btnRemoveChat),
	)

	return
}

func CreateImportedChats(chats []model.Chat) (menu *telebot.ReplyMarkup, err error) {
	if len(chats) <= 0 {
		return nil, fmt.Errorf("invalid amount of chats")
	}

	menu = &telebot.ReplyMarkup{}

	var chatButtons []telebot.Btn
	for _, chat := range chats {
		btn := menu.Data(chat.Name, fmt.Sprintf("chat_id_%s", chat.Id))
		chatButtons = append(chatButtons, btn)
	}

	rows := make([]telebot.Row, len(chatButtons))
	for i, btn := range chatButtons {
		rows[i] = menu.Row(btn)
	}

	menu.Inline(rows...)

	return
}
