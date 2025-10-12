package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/sharon-xa/gomux/tui/screens"
	"github.com/sharon-xa/gomux/tui/shared"
)

type Model struct {
	currentScreen shared.ScreenState
	state         shared.AppState

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
	state := shared.NewDefaultAppState()

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}

	return Model{
		currentScreen: shared.MAIN_MENU,
		state:         state,
		mainMenu:      screens.NewMainMenuModel(width, height),
		createScript:  screens.NewCreateScriptModel(width, height),
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
		m.mainMenu = m.mainMenu.SetSize(msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case shared.ScreenChangeMsg:
		m.currentScreen = msg.Screen

		switch msg.Screen {
		// case shared.EDIT_SCRIPT:
		// 	if scriptID, ok := msg.Data.(string); ok {
		// 		m.state.CurrentScriptID = scriptID
		// 		if script := m.state.GetScriptByID(scriptID); script != nil {
		// 			m.editScript = screens.NewEditScriptModelWithScript(*script)
		// 		}
		// 	}
		case shared.RUN_SCRIPT:
			if scriptID, ok := msg.Data.(string); ok {
				m.state.CurrentScriptID = scriptID
				if script := m.state.GetScriptByID(scriptID); script != nil {
					m.runScript = screens.NewRunScriptModelWithScript(*script)
				}
			}

		case shared.LIST_SCRIPTS:
			m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)

		case shared.MAIN_MENU:
			m.state.CurrentScriptID = ""
		}

		return m, nil

	case shared.ScriptCreatedMsg:
		m.state.AddScript(msg.Script)
		m.currentScreen = shared.LIST_SCRIPTS
		m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
		return m, nil

	// case shared.ScriptUpdatedMsg:
	// 	// Update script in state
	// 	m.state.UpdateScript(msg.ScriptID, msg.Script)
	// 	// Navigate back to list
	// 	m.currentScreen = shared.LIST_SCRIPTS
	// 	m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
	// 	return m, nil
	//
	// case shared.ScriptDeletedMsg:
	// 	// Delete script from state
	// 	m.state.DeleteScript(msg.ScriptID)
	// 	// Refresh list
	// 	m.listScripts = screens.NewListScriptsModelWithScripts(m.state.Scripts)
	// 	return m, nil

	case shared.ScriptExecutionCompleteMsg:
		m.state.AddExecutionResult(msg.Result)
		return m, nil

	case shared.SettingsUpdatedMsg:
		m.state.Settings = msg.Settings
		return m, nil

	case shared.BackMsg:
		m.currentScreen = shared.MAIN_MENU
		return m, nil

	case shared.ErrorMsg:
		m.state.LastError = msg.Error()
		return m, nil
	}

	switch m.currentScreen {
	case shared.MAIN_MENU:
		m.mainMenu, cmd = m.mainMenu.Update(msg)

	case shared.CREATE_SCRIPT:
		m.createScript, cmd = m.createScript.Update(msg)

	// case shared.EDIT_SCRIPT:
	// 	m.editScript, cmd = m.editScript.Update(msg)
	//
	// case shared.RUN_SCRIPT:
	// 	m.runScript, cmd = m.runScript.Update(msg)

	case shared.LIST_SCRIPTS:
		m.listScripts, cmd = m.listScripts.Update(msg)

		// case shared.SETTINGS:
		// 	m.settings, cmd = m.settings.Update(msg)
		//
		// case shared.EXECUTION_HISTORY:
		// 	m.executionHistory, cmd = m.executionHistory.Update(msg)
	}

	return m, cmd
}

// View renders the current screen
func (m Model) View() string {
	// Delegate to current screen
	switch m.currentScreen {
	case shared.MAIN_MENU:
		return m.mainMenu.View()

	case shared.CREATE_SCRIPT:
		return m.createScript.View()

	// case shared.EDIT_SCRIPT:
	// 	return m.editScript.View()

	case shared.RUN_SCRIPT:
		return m.runScript.View()

	case shared.LIST_SCRIPTS:
		return m.listScripts.View()

	// case shared.SETTINGS:
	// 	return m.settings.View()
	//
	// case shared.EXECUTION_HISTORY:
	// 	return m.executionHistory.View()

	default:
		return "Unknown screen"
	}
}

func (m *Model) GetScripts() []shared.Script {
	return m.state.Scripts
}

func (m *Model) GetSettings() shared.Settings {
	return m.state.Settings
}

func (m *Model) GetExecutionHistory() []shared.ScriptExecutionResult {
	return m.state.ExecutionHistory
}
