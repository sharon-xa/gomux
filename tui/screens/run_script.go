package screens

import (
	"github.com/sharon-xa/gomux/internal/models"
)

type RunScriptModel struct {
	scripts []models.Script

	cursor int
	width  int
	height int
}

func NewRunScriptModel(width, height int) RunScriptModel {
	return RunScriptModel{
		scripts: []models.Script{},
		cursor:  0,
	}
}

func NewRunScriptModelWithScript(script models.Script) RunScriptModel {
	return RunScriptModel{
		scripts: []models.Script{
			script,
		},
		cursor: 1,
	}
}

func (m RunScriptModel) View() string {
	return ""
}
