package shared

import "github.com/sharon-xa/gomux/internal/models"

type ScreenChangeMsg struct {
	Screen ScreenState
	Data   any
}

func NewScreenChange(screen ScreenState, data any) ScreenChangeMsg {
	return ScreenChangeMsg{
		Screen: screen,
		Data:   data,
	}
}

type ScriptCreatedMsg struct {
	Script models.Script
}

type ScriptUpdatedMsg struct {
	ScriptID string
	Script   models.Script
}

type ScriptDeletedMsg struct {
	ScriptID string
}

type ScriptSelectedMsg struct {
	ScriptID string
}

type ScriptExecutionStartMsg struct {
	ScriptID string
}

type ScriptExecutionCompleteMsg struct {
	Result ScriptExecutionResult
}

type SettingsUpdatedMsg struct {
	Settings Settings
}

// SearchQueryMsg is sent when search query changes
type SearchQueryMsg struct {
	Query string
}

// FilterTagsMsg is sent when filter tags change
type FilterTagsMsg struct {
	Tags []string
}

// ConfirmActionMsg is sent to confirm a destructive action
type ConfirmActionMsg struct {
	Action    string // "delete", "overwrite", etc.
	Message   string // Confirmation message to display
	OnConfirm any    // Data to pass if confirmed
	OnCancel  any    // Data to pass if cancelled
}

func NewConfirmation(action, message string, onConfirm, onCancel any) ConfirmActionMsg {
	return ConfirmActionMsg{
		Action:    action,
		Message:   message,
		OnConfirm: onConfirm,
		OnCancel:  onCancel,
	}
}

// CancelActionMsg is sent to cancel the current operation
type CancelActionMsg struct{}

// BackMsg is sent to go back to the previous screen
type BackMsg struct{}

// RefreshMsg is sent to refresh the current screen data
type RefreshMsg struct{}

// LoadingMsg is sent to indicate loading state
type LoadingMsg struct {
	IsLoading bool
	Message   string
}
