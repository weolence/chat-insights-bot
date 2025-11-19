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
	StateFileOfChatAwaiting  // bot asks user for file

	StateChatSelected // user selected chat

	StateChatAnalyzationSelected             // user selected analyzation of chat
	StateDateBoundsForSummarizingAwaiting    // bot asks user for date bounds of summarizing
	StateDescriptionForContextSearchAwaiting // bot asks user for description of context

	StateChatSettingsSelected   // user selected settings of chat
	StateNewNameForChatAwaiting // bot asks user for new chat name
)

// structure for keeping current user state
type UserData struct {
	State          UserState
	SelectedChatId int64
	NewChatName    string
}

// map of user states: map[telegramUserId]UserData
var UserStates sync.Map
