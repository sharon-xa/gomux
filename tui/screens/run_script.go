package screens

import "github.com/sharon-xa/gomux/tui/shared"

type RunScriptModel struct {
	scripts []shared.Script

	cursor int
	width  int
	height int
}

func NewRunScriptModel(width, height int) RunScriptModel {
	return RunScriptModel{
		scripts: []shared.Script{},
		cursor:  0,
	}
}

func NewRunScriptModelWithScript(script shared.Script) RunScriptModel {
	return RunScriptModel{
		scripts: []shared.Script{
			script,
		},
		cursor: 1,
	}
}

func (m RunScriptModel) View() string {
	return ""
}
