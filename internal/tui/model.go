package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/sharon-xa/gomux/internal/tui/messages"
	"github.com/sharon-xa/gomux/internal/tui/screens"
	"github.com/sharon-xa/gomux/internal/tui/style"
	"github.com/sharon-xa/gomux/internal/workspace"
)

type Model struct {
	currentScreen messages.ScreenState
	state         AppState

	mainMenu     screens.MainMenuModel
	createScript screens.CreateScriptModel
	listScripts  screens.ListScriptsModel
	runScript    screens.RunScriptModel
	// editScript       screens.EditScriptModel
	// settings         screens.SettingsModel
	// executionHistory screens.ExecutionHistoryModel

	width  int
	height int
}

func NewModel() Model {
	state := NewDefaultAppState()

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}

	return Model{
		currentScreen: messages.MAIN_MENU,
		state:         state,
		mainMenu:      screens.NewMainMenuModel(width, height),
		createScript:  screens.NewCreateScriptModel(),
		listScripts:   screens.NewListScriptsModel(width, height),
		runScript:     screens.NewRunScriptModel(width, height),
		width:         width,
		height:        height,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.PageStyle.GetFrameSize()
		m.mainMenu = m.mainMenu.SetSize(msg.Width-h, msg.Height-v)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case messages.ScreenChangeMsg:
		m.currentScreen = msg.Screen

		switch msg.Screen {
		// case messages.EDIT_SCRIPT:
		// 	if scriptID, ok := msg.Data.(string); ok {
		// 		m.state.CurrentScriptID = scriptID
		// 		if script := m.state.GetScriptByID(scriptID); script != nil {
		// 			m.editScript = screens.NewEditScriptModelWithScript(*script)
		// 		}
		// 	}
		case messages.RUN_SCRIPT:
			if scriptID, ok := msg.Data.(string); ok {
				m.state.CurrentScriptID = scriptID
				if script := m.state.GetScriptByID(scriptID); script != nil {
					m.runScript = screens.NewRunScriptModelWithScript(*script)
				}
			}

		case messages.LIST_SCRIPTS:
			m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)

		case messages.MAIN_MENU:
			m.state.CurrentScriptID = ""
		}

		return m, nil

	case messages.ScriptCreatedMsg:
		m.state.AddScript(msg.Script)
		m.currentScreen = messages.LIST_SCRIPTS
		m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
		return m, nil

	// case messages.ScriptUpdatedMsg:
	// 	// Update script in state
	// 	m.state.UpdateScript(msg.ScriptID, msg.Script)
	// 	// Navigate back to list
	// 	m.currentScreen = messages.LIST_SCRIPTS
	// 	m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
	// 	return m, nil
	//
	// case messages.ScriptDeletedMsg:
	// 	// Delete script from state
	// 	m.state.DeleteScript(msg.ScriptID)
	// 	// Refresh list
	// 	m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
	// 	return m, nil

	case messages.ScriptExecutionCompleteMsg:
		m.state.AddExecutionResult(msg.Result)
		return m, nil

	case messages.SettingsUpdatedMsg:
		m.state.Settings = msg.Settings
		return m, nil

	case messages.BackMsg:
		m.currentScreen = messages.MAIN_MENU
		return m, nil

	case messages.ErrorMsg:
		m.state.LastError = msg.Error()
		return m, nil
	}

	switch m.currentScreen {
	case messages.MAIN_MENU:
		m.mainMenu, cmd = m.mainMenu.Update(msg)

	case messages.CREATE_SCRIPT:
		m.createScript, cmd = m.createScript.Update(msg)

	// case messages.EDIT_SCRIPT:
	// 	m.editScript, cmd = m.editScript.Update(msg)
	//
	// case messages.RUN_SCRIPT:
	// 	m.runScript, cmd = m.runScript.Update(msg)

	case messages.LIST_SCRIPTS:
		m.listScripts, cmd = m.listScripts.Update(msg)

		// case messages.SETTINGS:
		// 	m.settings, cmd = m.settings.Update(msg)
	}

	return m, cmd
}

// View renders the current screen
func (m Model) View() string {
	// Delegate to current screen
	switch m.currentScreen {
	case messages.MAIN_MENU:
		return m.mainMenu.View()

	case messages.CREATE_SCRIPT:
		return m.createScript.View()

	// case messages.EDIT_SCRIPT:
	// 	return m.editScript.View()

	case messages.RUN_SCRIPT:
		return m.runScript.View()

	case messages.LIST_SCRIPTS:
		return m.listScripts.View()

	// case messages.SETTINGS:
	// 	return m.settings.View()
	//
	// case messages.EXECUTION_HISTORY:
	// 	return m.executionHistory.View()

	default:
		return "Unknown screen"
	}
}

func (m *Model) GetScripts() []workspace.Script {
	return m.state.Scripts
}

func (m *Model) GetSettings() messages.Settings {
	return m.state.Settings
}

func (m *Model) GetExecutionHistory() []messages.ScriptExecutionResult {
	return m.state.ExecutionHistory
}
