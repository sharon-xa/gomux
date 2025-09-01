package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	state          sessionState
	list           list.Model
	width          int
	height         int
	selectedScript string
	scripts        []string // Available scripts
	err            error
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)
)

func InitialModel() model {
	// Create menu items
	items := []list.Item{
		menuItem{
			title:       "📝 Create Script",
			description: "Create a new tmux script with guided setup",
			action:      createScript,
		},
		menuItem{
			title:       "✏️  Edit Script",
			description: "Edit an existing tmux script",
			action:      editScript,
		},
		menuItem{
			title:       "▶️  Run Script",
			description: "Run an existing tmux script",
			action:      runScript,
		},
		menuItem{
			title:       "📋 List Scripts",
			description: "View all available tmux scripts",
			action:      listScripts,
		},
		menuItem{
			title:       "⚙️  Settings",
			description: "Configure Gomux settings",
			action:      settings,
		},
	}

	// Create list with custom styling
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Gomux - Tmux Script Manager"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		state:   mainMenu,
		list:    l,
		scripts: []string{"dev-setup", "docker-env", "monitoring"}, // Mock data
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
