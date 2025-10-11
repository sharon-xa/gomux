package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	u "github.com/sharon-xa/gomux/tui/utils"
)

func (m model) View() string {
	switch m.state {
	case mainMenu:
		return m.viewMainMenu()
	case createScript:
		return m.viewCreateScript()
	case editScript:
		return m.viewEditScript()
	case runScript:
		return m.viewRunScript()
	case listScripts:
		return m.viewListScripts()
	case settings:
		return m.viewSettings()
	default:
		return "Unknown state"
	}
}

func (m model) viewMainMenu() string {
	return m.list.View()
}

func (m model) viewCreateScript() string {
	var s strings.Builder

	s.WriteString(u.TitleStyle.Render("📝 Create New Tmux Script"))
	s.WriteString("\n\n")

	s.WriteString("This will launch the interactive script creation wizard.\n\n")

	s.WriteString("The wizard will guide you through:\n")
	s.WriteString("• Session configuration (name, directory, type)\n")
	s.WriteString("• Window setup (commands, layouts, sync panes)\n")
	s.WriteString("• Pane configuration (splits, commands)\n")
	s.WriteString("• Review and save\n\n")

	s.WriteString(u.SuccessStyle.Render("Press Enter to start the wizard"))
	s.WriteString("\n\n")
	s.WriteString(u.HelpStyle.Render("Enter: Start wizard • Esc: Back to main menu"))

	return s.String()
}

func (m model) viewEditScript() string {
	var s strings.Builder

	s.WriteString(u.TitleStyle.Render("✏️ Edit Tmux Script"))
	s.WriteString("\n\n")

	s.WriteString("Select a script to edit:\n\n")

	for i, script := range m.scripts {
		s.WriteString(fmt.Sprintf("  %d. %s\n", i+1, script))
	}

	s.WriteString("\n")
	if m.selectedScript != "" {
		s.WriteString(u.SuccessStyle.Render(fmt.Sprintf("Selected: %s", m.selectedScript)))
		s.WriteString("\n")
	}

	s.WriteString(u.HelpStyle.Render("1-3: Select script • Esc: Back to main menu"))

	return s.String()
}

func (m model) viewRunScript() string {
	var s strings.Builder

	s.WriteString(u.TitleStyle.Render("▶️ Run Tmux Script"))
	s.WriteString("\n\n")

	s.WriteString("Select a script to run:\n\n")

	for i, script := range m.scripts {
		status := "✅ Ready"
		if i == 1 {
			status = "🔄 Running"
		}
		s.WriteString(fmt.Sprintf("  %d. %-15s %s\n", i+1, script, status))
	}

	s.WriteString("\n")
	if m.selectedScript != "" {
		s.WriteString(u.SuccessStyle.Render(fmt.Sprintf("Running: %s", m.selectedScript)))
		s.WriteString("\n")
	}

	s.WriteString(u.HelpStyle.Render("1-3: Run script • Esc: Back to main menu"))

	return s.String()
}

func (m model) viewListScripts() string {
	var s strings.Builder

	s.WriteString(u.TitleStyle.Render("📋 Tmux Scripts"))
	s.WriteString("\n\n")

	if len(m.scripts) == 0 {
		s.WriteString("No scripts found.\n\n")
		s.WriteString("Create your first script from the main menu!")
	} else {
		s.WriteString(fmt.Sprintf("Found %d script(s):\n\n", len(m.scripts)))

		// Table header
		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
		s.WriteString(headerStyle.Render("Name"))
		s.WriteString(strings.Repeat(" ", 15))
		s.WriteString(headerStyle.Render("Status"))
		s.WriteString(strings.Repeat(" ", 10))
		s.WriteString(headerStyle.Render("Last Modified"))
		s.WriteString("\n")
		s.WriteString(strings.Repeat("─", 50))
		s.WriteString("\n")

		for i, script := range m.scripts {
			status := "✅ Ready"
			lastMod := "2024-01-15"
			if i == 1 {
				status = "🔄 Running"
				lastMod = "2024-01-20"
			}

			s.WriteString(fmt.Sprintf("%-20s %-15s %s\n", script, status, lastMod))
		}
	}

	s.WriteString("\n")
	s.WriteString(u.HelpStyle.Render("d: Delete script • Esc: Back to main menu"))

	return s.String()
}

func (m model) viewSettings() string {
	var s strings.Builder

	s.WriteString(u.TitleStyle.Render("⚙️ Gomux Settings"))
	s.WriteString("\n\n")

	s.WriteString("Configuration:\n\n")

	settings := []struct {
		name  string
		value string
	}{
		{"Scripts Directory", "~/.config/gomux/scripts"},
		{"Default Shell", "/bin/zsh"},
		{"Auto-start Sessions", "enabled"},
		{"Theme", "default"},
		{"Editor", "vim"},
	}

	for _, setting := range settings {
		s.WriteString(fmt.Sprintf("  %-20s: %s\n", setting.name, setting.value))
	}

	s.WriteString("\n")
	s.WriteString(
		u.HelpStyle.Render("e: Edit settings • r: Reset to defaults • Esc: Back to main menu"),
	)

	return s.String()
}
