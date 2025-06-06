package tui

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	count int
}

func InitialModel() model {
	return model{count: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}
