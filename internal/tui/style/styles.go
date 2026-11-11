package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PageStyle = lipgloss.NewStyle().
			Margin(2, 2, 2, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Bold(true)

	FocusedLabelStyle = LabelStyle.Foreground(
		lipgloss.Color("#7D56F4"),
	)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)
)

func StyledPage(page string) string {
	return PageStyle.Render(page)
}

func StyledTitle(title string) string {
	return TitleStyle.Render(title)
}
