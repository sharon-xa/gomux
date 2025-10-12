package shared

type ScreenState int

const (
	MAIN_MENU ScreenState = iota
	// Main Screens
	CREATE_SCRIPT
	EDIT_SCRIPT
	RUN_SCRIPT
	LIST_SCRIPTS
	SETTINGS
	EXECUTION_HISTORY
)

func (s ScreenState) String() string {
	switch s {
	case MAIN_MENU:
		return "Main Menu"
	case CREATE_SCRIPT:
		return "Create Script"
	case EDIT_SCRIPT:
		return "Edit Script"
	case RUN_SCRIPT:
		return "Run Script"
	case LIST_SCRIPTS:
		return "List Scripts"
	case SETTINGS:
		return "Settings"
	case EXECUTION_HISTORY:
		return "Execution History"
	default:
		return "Unknown"
	}
}
