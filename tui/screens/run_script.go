package screens

import (
	"github.com/sharon-xa/gomux/internal/workspace"
)

type RunScriptModel struct {
	scripts []workspace.Script

	cursor int
	width  int
	height int
}

func NewRunScriptModel(width, height int) RunScriptModel {
	return RunScriptModel{
		scripts: []workspace.Script{},
		cursor:  0,
	}
}

func NewRunScriptModelWithScript(script workspace.Script) RunScriptModel {
	return RunScriptModel{
		scripts: []workspace.Script{
			script,
		},
		cursor: 1,
	}
}

func (m RunScriptModel) View() string {
	return ""
}
