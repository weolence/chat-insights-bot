package model

import "sync"

// structure for fetching already created users from database
type User struct {
	TelegramId int64 `json:"telegram_id"`
}

type UserState int

const (
	StateNone UserState = iota

	StateRootMenu // user stands in root menu

	StateChatImportSelected  // user pressed import chat
	StateNameForChatAwaiting // bot asks user for chat name
	StateNameForChatReceived // bot received name for chat from user
	StateFileOfChatAwaiting  // bot asks user for file
	StateFileOfChatReceived  // bot received file of chat from user

	StateChatSelected // user selected chat

	StateChatAnalyzationSelected             // user selected analyzation of chat
	StateDateBoundsForSummarizingAwaiting    // bot asks user for date bounds of summarizing
	StateDateBoundsForSummarizingReceived    // bot received date bounds from user
	StateDescriptionForContextSearchAwaiting // bot asks user for description of context
	StateDescriptionForContextSearchReceived // bot received description of context from user

	StateChatSettingsSelected   // user selected settings of chat
	StateNewNameForChatAwaiting // bot asks user for new chat name
	StateNewNameForChatReceived // bot received new name for chat from user
)

// structure for keeping current user state
type UserData struct {
	State          UserState
	SelectedChatId int64
	NewChatName    string
}

// map of user states: map[telegramUserId]UserData
var UserStates sync.Map
