package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sharon-xa/gomux/internal/workspace"
	"github.com/sharon-xa/gomux/pkg/actions"
	"github.com/sharon-xa/gomux/tui/commenui"
	"github.com/sharon-xa/gomux/tui/shared"
)

// Note: This only stores UI state, NOT the actual scripts list
type CreateScriptModel struct {
	// Form fields (local UI state only)
	nameInput        textinput.Model
	descriptionInput textinput.Model
	contentInput     textarea.Model
	tagsInput        textinput.Model
	languageInput    textinput.Model
	actions          actions.IActions

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
func NewCreateScriptModel() CreateScriptModel {
	// name input
	nameInput := textinput.New()
	nameInput.Placeholder = "My Script"
	nameInput.Focus()
	nameInput.CharLimit = 100

	// description input
	descInput := textinput.New()
	descInput.Placeholder = "What does this script do?"
	descInput.CharLimit = 200

	// language input
	langInput := textinput.New()
	langInput.Placeholder = "bash"
	langInput.CharLimit = 20

	// tags input
	tagsInput := textinput.New()
	tagsInput.Placeholder = "tag1, tag2, tag3"
	tagsInput.CharLimit = 100

	// content textarea
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
		actions:          actions.NewActions(),
	}
}

// Init initializes the create script screen
func (m CreateScriptModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the create script screen
func (m CreateScriptModel) Update(msg tea.Msg) (CreateScriptModel, tea.Cmd) {
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

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.focusedField {
	case nameField:
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)

	case descriptionField:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		cmds = append(cmds, cmd)

	case languageField:
		m.languageInput, cmd = m.languageInput.Update(msg)
		cmds = append(cmds, cmd)

	case tagsField:
		m.tagsInput, cmd = m.tagsInput.Update(msg)
		cmds = append(cmds, cmd)

	case contentField:
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

	session := workspace.Session{}
	window := workspace.Window{}
	startingWindow := 1
	pane := workspace.Pane{}

	script, err := m.actions.NewScript(
		session,
		[]workspace.Window{window},
		startingWindow,
		[]workspace.Pane{pane},
	)
	if err != nil {
		// handle error
	}

	// Return a command that sends ScriptCreatedMsg
	return func() tea.Msg {
		return shared.ScriptCreatedMsg{Script: *script}
	}
}

func (m CreateScriptModel) View() string {
	var s strings.Builder

	s.WriteString(commenui.StyledTitle("📝 Create New Script"))
	s.WriteString("\n\n")

	nameLabel := commenui.LabelStyle.Render("Name:")
	if m.focusedField == nameField {
		nameLabel = commenui.FocusedLabelStyle.Render("Name:")
	}
	s.WriteString(nameLabel + "\n")
	s.WriteString(m.nameInput.View() + "\n\n")

	descLabel := commenui.LabelStyle.Render("Description:")
	if m.focusedField == descriptionField {
		descLabel = commenui.FocusedLabelStyle.Render("Description:")
	}
	s.WriteString(descLabel + "\n")
	s.WriteString(m.descriptionInput.View() + "\n\n")

	langLabel := commenui.LabelStyle.Render("Language:")
	if m.focusedField == languageField {
		langLabel = commenui.FocusedLabelStyle.Render("Language:")
	}
	s.WriteString(langLabel + "\n")
	s.WriteString(m.languageInput.View() + "\n\n")

	tagsLabel := commenui.LabelStyle.Render("Tags (comma-separated):")
	if m.focusedField == tagsField {
		tagsLabel = commenui.FocusedLabelStyle.Render("Tags (comma-separated):")
	}
	s.WriteString(tagsLabel + "\n")
	s.WriteString(m.tagsInput.View() + "\n\n")

	contentLabel := commenui.LabelStyle.Render("Content:")
	if m.focusedField == contentField {
		contentLabel = commenui.FocusedLabelStyle.Render("Content:")
	}
	s.WriteString(contentLabel + "\n")
	s.WriteString(m.contentInput.View() + "\n\n")

	if m.err != nil {
		s.WriteString(commenui.ErrorStyle.Render("❌ "+m.err.Error()) + "\n\n")
	}

	helpMsg := commenui.GetHelpMessage([]commenui.HelpInstructions{
		{Keys: []string{"tab", "shift+tab"}, Action: "navigate"},
		commenui.SaveInstruction,
		commenui.CancelInstruction,
	})
	s.WriteString(helpMsg)

	return commenui.StyledPage(s.String())
}
