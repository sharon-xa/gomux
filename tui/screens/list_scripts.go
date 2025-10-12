package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sharon-xa/gomux/tui/shared"
)

// ListScriptsModel represents the list scripts screen
// Note: It holds a COPY of the scripts for display, not the source of truth
type ListScriptsModel struct {
	// Local state for UI only
	scripts       []shared.Script // Read-only copy from root model
	cursor        int             // Selected script index
	width         int
	height        int
	showingDetail bool // Whether showing detail view
}

// NewListScriptsModel creates a new list scripts model
func NewListScriptsModel(width, height int) ListScriptsModel {
	return ListScriptsModel{
		scripts:       []shared.Script{},
		cursor:        0,
		showingDetail: false,
	}
}

// NewListScriptsModelWithScripts creates model with scripts
// This is called by Root Model when transitioning to this screen
func NewListScriptsModelWithScripts(scripts []shared.Script) ListScriptsModel {
	return ListScriptsModel{
		scripts:       scripts,
		cursor:        0,
		showingDetail: false,
	}
}

// Init initializes the list scripts screen
func (m ListScriptsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the list scripts screen
func (m ListScriptsModel) Update(msg tea.Msg) (ListScriptsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// If in detail view, different keys apply
		if m.showingDetail {
			switch msg.String() {
			case "esc", "q":
				m.showingDetail = false
				return m, nil
			case "e":
				// Edit selected script
				if len(m.scripts) > 0 {
					scriptID := m.scripts[m.cursor].ID
					return m, func() tea.Msg {
						return shared.NewScreenChange(shared.EDIT_SCRIPT, scriptID)
					}
				}
			case "r":
				// Run selected script
				if len(m.scripts) > 0 {
					scriptID := m.scripts[m.cursor].ID
					return m, func() tea.Msg {
						return shared.NewScreenChange(shared.RUN_SCRIPT, scriptID)
					}
				}
			case "d":
				// Delete selected script
				if len(m.scripts) > 0 {
					scriptID := m.scripts[m.cursor].ID
					return m, func() tea.Msg {
						return shared.ScriptDeletedMsg{ScriptID: scriptID}
					}
				}
			}
			return m, nil
		}

		// Normal list view keys
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.scripts)-1 {
				m.cursor++
			}

		case "enter", " ":
			// Show detail view
			m.showingDetail = true

		case "e":
			// Quick edit without showing detail
			if len(m.scripts) > 0 {
				scriptID := m.scripts[m.cursor].ID
				return m, func() tea.Msg {
					return shared.NewScreenChange(shared.EDIT_SCRIPT, scriptID)
				}
			}

		case "r":
			// Quick run without showing detail
			if len(m.scripts) > 0 {
				scriptID := m.scripts[m.cursor].ID
				return m, func() tea.Msg {
					return shared.NewScreenChange(shared.RUN_SCRIPT, scriptID)
				}
			}

		case "n":
			// Create new script
			return m, func() tea.Msg {
				return shared.NewScreenChange(shared.CREATE_SCRIPT, nil)
			}

		case "d":
			// Delete script
			if len(m.scripts) > 0 {
				scriptID := m.scripts[m.cursor].ID
				return m, func() tea.Msg {
					return shared.ScriptDeletedMsg{ScriptID: scriptID}
				}
			}

		case "esc", "q":
			// Go back to main menu
			return m, func() tea.Msg {
				return shared.BackMsg{}
			}

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the list scripts screen
func (m ListScriptsModel) View() string {
	if m.showingDetail {
		return m.renderDetailView()
	}
	return m.renderListView()
}

// renderListView renders the list of scripts
func (m ListScriptsModel) renderListView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ADD8")).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#00ADD8")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(1)

	var s strings.Builder

	// Title
	title := fmt.Sprintf("📋 Scripts (%d)", len(m.scripts))
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	// If no scripts
	if len(m.scripts) == 0 {
		s.WriteString("  No scripts yet. Press 'n' to create one!\n")
	} else {
		// List scripts
		for i, script := range m.scripts {
			cursor := "  "
			if m.cursor == i {
				cursor = "▶ "
			}

			// Format script info
			lastRun := "never"
			if script.LastRunAt != nil {
				lastRun = formatTime(*script.LastRunAt)
			}

			favorite := ""
			if script.Favorite {
				favorite = "⭐ "
			}

			line := fmt.Sprintf("%s%s [%s] - Last run: %s",
				favorite,
				script.Name,
				script.Language,
				lastRun,
			)

			// Apply style
			var rendered string
			if m.cursor == i {
				rendered = selectedStyle.Render(line)
			} else {
				rendered = normalStyle.Render(line)
			}

			s.WriteString(cursor + rendered + "\n")
		}
	}

	// Help text
	s.WriteString("\n")
	help := helpStyle.Render(
		"↑/k up • ↓/j down • enter details • e edit • r run • d delete • n new • esc back",
	)
	s.WriteString(help)

	return s.String()
}

// renderDetailView renders detailed view of selected script
func (m ListScriptsModel) renderDetailView() string {
	if len(m.scripts) == 0 || m.cursor >= len(m.scripts) {
		return "No script selected"
	}

	script := m.scripts[m.cursor]

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ADD8")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	codeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a1a")).
		Foreground(lipgloss.Color("#00FF00")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(1)

	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("📄 Script Details"))
	s.WriteString("\n\n")

	// Name
	s.WriteString(labelStyle.Render("Name: "))
	s.WriteString(contentStyle.Render(script.Name))
	s.WriteString("\n\n")

	// Description
	s.WriteString(labelStyle.Render("Description: "))
	desc := script.Description
	if desc == "" {
		desc = "(no description)"
	}
	s.WriteString(contentStyle.Render(desc))
	s.WriteString("\n\n")

	// Language
	s.WriteString(labelStyle.Render("Language: "))
	s.WriteString(contentStyle.Render(script.Language))
	s.WriteString("\n\n")

	// Tags
	s.WriteString(labelStyle.Render("Tags: "))
	if len(script.Tags) == 0 {
		s.WriteString(contentStyle.Render("(no tags)"))
	} else {
		s.WriteString(contentStyle.Render(strings.Join(script.Tags, ", ")))
	}
	s.WriteString("\n\n")

	// Created/Updated
	s.WriteString(labelStyle.Render("Created: "))
	s.WriteString(contentStyle.Render(formatTime(script.CreatedAt)))
	s.WriteString("\n")

	s.WriteString(labelStyle.Render("Updated: "))
	s.WriteString(contentStyle.Render(formatTime(script.UpdatedAt)))
	s.WriteString("\n")

	// Last Run
	if script.LastRunAt != nil {
		s.WriteString(labelStyle.Render("Last Run: "))
		s.WriteString(contentStyle.Render(formatTime(*script.LastRunAt)))
		s.WriteString("\n")
	}

	// Content preview
	s.WriteString("\n")
	s.WriteString(labelStyle.Render("Content:"))
	s.WriteString("\n")

	// Show first 10 lines of content
	lines := strings.Split(script.Content, "\n")
	preview := lines
	if len(lines) > 10 {
		preview = lines[:10]
	}
	s.WriteString(codeStyle.Render(strings.Join(preview, "\n")))

	if len(lines) > 10 {
		s.WriteString(contentStyle.Render(fmt.Sprintf("... (%d more lines)", len(lines)-10)))
		s.WriteString("\n")
	}

	// Help text
	s.WriteString("\n")
	help := helpStyle.Render("e edit • r run • d delete • esc back to list")
	s.WriteString(help)

	return s.String()
}

// formatTime formats a time.Time into a human-readable string
func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}
