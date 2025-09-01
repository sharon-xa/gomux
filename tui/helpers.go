package tui

import "strings"

func (m *model) updateSize() {
	m.list.SetWidth(m.width)
	m.list.SetHeight(m.height - 4) // Account for title and help text
}

func centerText(text string, width int) string {
	if width <= len(text) {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}
