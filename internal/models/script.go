package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Script struct {
	Session        Session  `json:"session"`
	Windows        []Window `json:"windows"`
	StartingWindow int      `json:"startingWindow"` // when the tmux start on which window should we start
	Panes          []Pane   `json:"panes"`
}

func NewScript(s Session, ws []Window, sw int, ps []Pane) *Script {
	return &Script{
		Session:        s,
		Windows:        ws,
		StartingWindow: sw,
		Panes:          ps,
	}
}

func (s *Script) Jsonify() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (s *Script) Stringify() string {
	scriptTemp := make([]string, SCRIPT_TEMPLATE_SECTIONS)
	scriptTemp[0] = shebang()
	scriptTemp[1] = s.setScriptVars()
	scriptTemp[2] = checkDirectoryExistence()
	scriptTemp[3] = checkIfSessionExists()
	scriptTemp[4] = s.Session.createSession()
	scriptTemp[5] = s.createWindows()
	scriptTemp[6] = s.createPanes()
	scriptTemp[7] = s.setStartingWindow()

	return strings.Join(scriptTemp, SECTION_SEPARATOR)
}

func (s *Script) setScriptVars() string {
	return fmt.Sprintf("%s=\"%s\"\n%s=\"%s\"",
		SESSION_NAME_VARIABLE,
		s.Session.Name,
		SESSION_DIR_VARIABLE,
		s.Session.Directory,
	)
}

func (s *Script) createWindows() string {
	windows := ""

	for i, window := range s.Windows {
		windows += fmt.Sprintf("# Window Number %d", i)
		windows += window.createWindow()
		windows += "\n"
	}

	return windows
}

func (s *Script) setStartingWindow() string {
	return fmt.Sprintf(`tmux select-window -t $%s:%d`, SESSION_NAME_VARIABLE, s.StartingWindow)
}

func (s *Script) createPanes() string {
	panes := ""

	for i, pane := range s.Panes {
		panes += fmt.Sprintf("## Pane Number %d", i)
		panes += pane.createPane()
		panes += "\n"
	}

	return panes
}
