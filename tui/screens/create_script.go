package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/sharon-xa/gomux/tui/shared"
)

// CreateScriptModel represents the create script screen
// Note: This only stores UI state, NOT the actual scripts list
type CreateScriptModel struct {
	// Form fields (local UI state only)
	nameInput        textinput.Model
	descriptionInput textinput.Model
	contentInput     textarea.Model
	tagsInput        textinput.Model
	languageInput    textinput.Model

	// UI state
	focusedField int
	width        int
	height       int
	err          error
}

const (
	nameField = iota
	descriptionField
	languageField
	tagsField
	contentField
)

// NewCreateScriptModel creates a new create script model
func NewCreateScriptModel(width, height int) CreateScriptModel {
	// Initialize name input
	nameInput := textinput.New()
	nameInput.Placeholder = "My Script"
	nameInput.Focus()
	nameInput.CharLimit = 100

	// Initialize description input
	descInput := textinput.New()
	descInput.Placeholder = "What does this script do?"
	descInput.CharLimit = 200

	// Initialize language input
	langInput := textinput.New()
	langInput.Placeholder = "bash"
	langInput.CharLimit = 20

	// Initialize tags input
	tagsInput := textinput.New()
	tagsInput.Placeholder = "tag1, tag2, tag3"
	tagsInput.CharLimit = 100

	// Initialize content textarea
	contentInput := textarea.New()
	contentInput.Placeholder = "#!/bin/bash\necho 'Hello World'"
	contentInput.SetHeight(10)

	return CreateScriptModel{
		nameInput:        nameInput,
		descriptionInput: descInput,
		languageInput:    langInput,
		tagsInput:        tagsInput,
		contentInput:     contentInput,
		focusedField:     nameField,
	}
}

// Init initializes the create script screen
func (m CreateScriptModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the create script screen
func (m CreateScriptModel) Update(msg tea.Msg) (CreateScriptModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to main menu without saving
			return m, func() tea.Msg {
				return shared.BackMsg{}
			}

		case "ctrl+s":
			// Save the script
			return m, m.saveScript()

		case "tab", "shift+tab":
			// Navigate between fields
			m.nextField(msg.String() == "shift+tab")
			return m, nil

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// Update the focused field
	switch m.focusedField {
	case nameField:
		var cmd tea.Cmd
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)

	case descriptionField:
		var cmd tea.Cmd
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		cmds = append(cmds, cmd)

	case languageField:
		var cmd tea.Cmd
		m.languageInput, cmd = m.languageInput.Update(msg)
		cmds = append(cmds, cmd)

	case tagsField:
		var cmd tea.Cmd
		m.tagsInput, cmd = m.tagsInput.Update(msg)
		cmds = append(cmds, cmd)

	case contentField:
		var cmd tea.Cmd
		m.contentInput, cmd = m.contentInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// nextField moves focus to the next/previous field
func (m *CreateScriptModel) nextField(reverse bool) {
	if reverse {
		m.focusedField--
		if m.focusedField < 0 {
			m.focusedField = contentField
		}
	} else {
		m.focusedField++
		if m.focusedField > contentField {
			m.focusedField = nameField
		}
	}

	// Update focus state
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.languageInput.Blur()
	m.tagsInput.Blur()
	m.contentInput.Blur()

	switch m.focusedField {
	case nameField:
		m.nameInput.Focus()
	case descriptionField:
		m.descriptionInput.Focus()
	case languageField:
		m.languageInput.Focus()
	case tagsField:
		m.tagsInput.Focus()
	case contentField:
		m.contentInput.Focus()
	}
}

// saveScript creates a ScriptCreatedMsg with the form data
// This is the key: we don't modify shared data directly, we send a message
func (m CreateScriptModel) saveScript() tea.Cmd {
	// Validate
	if strings.TrimSpace(m.nameInput.Value()) == "" {
		m.err = fmt.Errorf("name is required")
		return nil
	}

	if strings.TrimSpace(m.contentInput.Value()) == "" {
		m.err = fmt.Errorf("content is required")
		return nil
	}

	// Parse tags
	tagsStr := m.tagsInput.Value()
	var tags []string
	if tagsStr != "" {
		for tag := range strings.SplitSeq(tagsStr, ",") {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	// Set default language if not provided
	language := m.languageInput.Value()
	if language == "" {
		language = "bash"
	}

	// Create the script object
	script := shared.Script{
		ID:          uuid.New().String(),
		Name:        strings.TrimSpace(m.nameInput.Value()),
		Description: strings.TrimSpace(m.descriptionInput.Value()),
		Content:     m.contentInput.Value(),
		Language:    language,
		Tags:        tags,
		Favorite:    false,
	}

	// Return a command that sends ScriptCreatedMsg
	// The Root Model will handle adding it to the shared state
	return func() tea.Msg {
		return shared.ScriptCreatedMsg{Script: script}
	}
}

// View renders the create script screen
func (m CreateScriptModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ADD8")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	focusedLabelStyle := labelStyle.Foreground(
		lipgloss.Color("#00ADD8"),
	)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("📝 Create New Script"))
	s.WriteString("\n\n")

	// Name field
	nameLabel := labelStyle.Render("Name:")
	if m.focusedField == nameField {
		nameLabel = focusedLabelStyle.Render("Name:")
	}
	s.WriteString(nameLabel + "\n")
	s.WriteString(m.nameInput.View() + "\n\n")

	// Description field
	descLabel := labelStyle.Render("Description:")
	if m.focusedField == descriptionField {
		descLabel = focusedLabelStyle.Render("Description:")
	}
	s.WriteString(descLabel + "\n")
	s.WriteString(m.descriptionInput.View() + "\n\n")

	// Language field
	langLabel := labelStyle.Render("Language:")
	if m.focusedField == languageField {
		langLabel = focusedLabelStyle.Render("Language:")
	}
	s.WriteString(langLabel + "\n")
	s.WriteString(m.languageInput.View() + "\n\n")

	// Tags field
	tagsLabel := labelStyle.Render("Tags (comma-separated):")
	if m.focusedField == tagsField {
		tagsLabel = focusedLabelStyle.Render("Tags (comma-separated):")
	}
	s.WriteString(tagsLabel + "\n")
	s.WriteString(m.tagsInput.View() + "\n\n")

	// Content field
	contentLabel := labelStyle.Render("Content:")
	if m.focusedField == contentField {
		contentLabel = focusedLabelStyle.Render("Content:")
	}
	s.WriteString(contentLabel + "\n")
	s.WriteString(m.contentInput.View() + "\n\n")

	// Error message
	if m.err != nil {
		s.WriteString(errorStyle.Render("❌ "+m.err.Error()) + "\n\n")
	}

	// Help text
	help := helpStyle.Render("tab/shift+tab navigate • ctrl+s save • esc cancel")
	s.WriteString(help)

	return s.String()
}
