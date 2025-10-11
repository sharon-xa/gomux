package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/sharon-xa/gomux/internal/models"
	"github.com/sharon-xa/gomux/pkg/actions"
	"github.com/sharon-xa/gomux/tui/utils"
)

type sessionState int

const (
	mainMenu sessionState = iota
	createScript
	editScript
	runScript
	listScripts
	settings
)

type menuItem struct {
	title       string
	description string
	action      sessionState
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

type model struct {
	// States
	state sessionState

	// lists
	list list.Model

	// Huh form
	form *huh.Form

	// view
	width  int
	height int

	// backend
	actions        actions.IActions
	script         models.Script
	session        models.Session
	windows        []models.Window
	panes          []models.Pane
	selectedScript string
	scripts        []string // Available scripts

	// error handling
	err error
}

func InitialModel() model {
	list := getMainMenu()

	return model{
		state:   mainMenu,
		list:    list,
		scripts: []string{"dev-setup", "docker-env", "monitoring"}, // Mock data
		actions: actions.NewActions(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) getSelectedMenuItem() menuItem {
	if item, ok := m.list.SelectedItem().(menuItem); ok {
		return item
	}
	return menuItem{}
}

func getMainMenu() list.Model {
	items := []list.Item{
		menuItem{
			title:       "📝  Create Script",
			description: "Create a new tmux script with guided setup",
			action:      createScript,
		},
		menuItem{
			title:       "🖊️  Edit Script",
			description: "Edit an existing tmux script",
			action:      editScript,
		},
		menuItem{
			title:       "▶️  Run Script",
			description: "Run an existing tmux script",
			action:      runScript,
		},
		menuItem{
			title:       "📋  List Scripts",
			description: "View all available tmux scripts",
			action:      listScripts,
		},
		menuItem{
			title:       "⚙️  Settings",
			description: "Configure Gomux settings",
			action:      settings,
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	utils.StyleList(&l)
	return l
}
