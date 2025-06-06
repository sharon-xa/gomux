package tui

import (
	"fmt"
)

func (m model) View() string {
	return fmt.Sprintf(
		"Welcome to Gomux!\n\nCount: %d\n\n↑ to increase, ↓ to decrease, q to quit.",
		m.count,
	)
}
