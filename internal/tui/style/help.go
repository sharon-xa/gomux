package style

import "strings"

type HelpInstructions struct {
	Keys   []string
	Action string
}

var (
	QuitInstruction = HelpInstructions{Keys: []string{"q"}, Action: "quit"}

	UpInstuction     = HelpInstructions{Keys: []string{"↑", "k"}, Action: "up"}
	DownInstuction   = HelpInstructions{Keys: []string{"↓", "j"}, Action: "down"}
	LeftInstruction  = HelpInstructions{Keys: []string{"←", "h"}, Action: "move left"}
	RightInstruction = HelpInstructions{Keys: []string{"→", "l"}, Action: "move right"}

	SaveInstruction    = HelpInstructions{Keys: []string{"ctrl+s"}, Action: "save"}
	CancelInstruction  = HelpInstructions{Keys: []string{"esc"}, Action: "cancel"}
	SelectInstuction   = HelpInstructions{Keys: []string{"enter"}, Action: "select"}
	RefreshInstruction = HelpInstructions{Keys: []string{"r"}, Action: "refresh"}
)

// Return a styled help message
func GetHelpMessage(instructions []HelpInstructions) string {
	// Example help message: "↑/k up • ↓/j down • enter select • q quit"
	var b strings.Builder
	for i, ins := range instructions {
		for j, key := range ins.Keys {
			b.WriteString(key)
			if j < len(ins.Keys)-1 {
				b.WriteString("/")
			}
		}

		b.WriteString(" " + ins.Action)

		if i < len(instructions)-1 {
			b.WriteString(" • ")
		}
	}

	return HelpStyle.Render(b.String())
}
