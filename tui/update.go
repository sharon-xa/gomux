package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSize()

	case tea.KeyMsg:
		switch m.state {
		case mainMenu:
			return m.updateMainMenu(msg)
		case createScript:
			return m.updateCreateScript(msg)
		case editScript:
			return m.updateEditScript(msg)
		case runScript:
			return m.updateRunScript(msg)
		case listScripts:
			return m.updateListScripts(msg)
		case settings:
			return m.updateSettings(msg)
		}
	}

	// Update the list if we're in main menu
	if m.state == mainMenu {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "enter":
		selectedItem := m.getSelectedMenuItem()
		m.state = selectedItem.action
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateCreateScript(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = mainMenu
		return m, nil
	case "enter":
		// Here you would launch the script creation flow
		// For now, just show a success message and return to menu
		m.state = mainMenu
		return m, nil
	}
	return m, nil
}

func (m model) updateEditScript(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = mainMenu
		return m, nil
	case "1", "2", "3":
		// Mock script selection
		scripts := []string{"dev-setup", "docker-env", "monitoring"}
		if idx := int(msg.String()[0] - '1'); idx < len(scripts) {
			m.selectedScript = scripts[idx]
			// Here you would launch the edit flow
			m.state = mainMenu
		}
		return m, nil
	}
	return m, nil
}

func (m model) updateRunScript(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = mainMenu
		return m, nil
	case "1", "2", "3":
		// Mock script running
		scripts := []string{"dev-setup", "docker-env", "monitoring"}
		if idx := int(msg.String()[0] - '1'); idx < len(scripts) {
			m.selectedScript = scripts[idx]
			// Here you would run the script
			m.state = mainMenu
		}
		return m, nil
	}
	return m, nil
}

func (m model) updateListScripts(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = mainMenu
		return m, nil
	case "d":
		// Delete selected script (mock)
		return m, nil
	}
	return m, nil
}

func (m model) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = mainMenu
		return m, nil
	}
	return m, nil
}
