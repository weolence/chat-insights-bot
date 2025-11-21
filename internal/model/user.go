package model

type UserState int

const (
	StateNone UserState = iota

	StateRootMenu // user stands in root menu

	StateNameForChatAwaiting // bot asks user for chat name
	StateFileOfChatAwaiting  // bot asks user for file

	StateChatSelectionAwaiting       // bot asks user for chat selecting
	StateChatInteractionTypeAwaiting // user selected chat

	StateChatAnalyzationTypeAwaiting         // user selected analyzation of chat
	StateDescriptionForContextSearchAwaiting // bot asks user for description of context

	StateChatSettingsTypeAwaiting // user selected settings of chat
	StateNewNameForChatAwaiting   // bot asks user for new chat name
)

// structure for fetching already created users from database
type User struct {
	TelegramId   int64 `json:"telegram_id"`
	State        UserState
	SelectedChat *Chat
	NewChatName  string
}
