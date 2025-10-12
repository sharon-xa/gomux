package screens

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sharon-xa/gomux/tui/commenui"
	"github.com/sharon-xa/gomux/tui/shared"
)

type Action int

const (
	createScript Action = iota
	editScript
	listScripts
	runScript
	settings
)

type menuItem struct {
	title       string
	description string
	action      Action
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

type MainMenuModel struct {
	list   list.Model
	width  int
	height int
}

func (m MainMenuModel) getSelectedMenuItem() menuItem {
	if item, ok := m.list.SelectedItem().(menuItem); ok {
		return item
	}
	return menuItem{}
}

func (m MainMenuModel) SetSize(width, height int) MainMenuModel {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
	return m
}

func NewMainMenuModel(width, height int) MainMenuModel {
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
			title:       "📋  List Scripts",
			description: "View all available tmux scripts",
			action:      listScripts,
		},
		menuItem{
			title:       "▶️  Run Script",
			description: "Run an existing tmux script",
			action:      runScript,
		},
		menuItem{
			title:       "⚙️  Settings",
			description: "Configure Gomux settings",
			action:      settings,
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), width, height)
	commenui.StyleList(&l)
	return MainMenuModel{
		list:   l,
		width:  width,
		height: height,
	}
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (MainMenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, m.handleSelection()

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) handleSelection() tea.Cmd {
	switch m.getSelectedMenuItem().action {
	case createScript:
		return func() tea.Msg {
			return shared.NewScreenChange(shared.CREATE_SCRIPT, nil)
		}

	case listScripts:
		return func() tea.Msg {
			return shared.NewScreenChange(shared.LIST_SCRIPTS, nil)
		}

	case runScript:
		return func() tea.Msg {
			return shared.NewScreenChange(shared.RUN_SCRIPT, nil)
		}

	case settings:
		return func() tea.Msg {
			return shared.NewScreenChange(shared.SETTINGS, nil)
		}

	default:
		return nil
	}
}

func (m MainMenuModel) View() string {
	log.Printf("List dimensions: width: %d - height: %d", m.list.Width(), m.list.Height())
	return m.list.View()
}
