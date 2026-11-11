// Package workspace holds the gomux spec: a tmux session with its windows and
// panes. It knows nothing about bash or storage; see internal/render for
// turning a spec into a script.
package workspace

import (
	"encoding/json"
)

type Script struct {
	Session        Session  `json:"session"`
	Windows        []Window `json:"windows"`
	StartingWindow int      `json:"startingWindow"` // when the tmux start on which window should we start
	Panes          []Pane   `json:"panes"`

	// TODO: maybe add some metadata like the created at and updated at timestamps
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

func LoadScriptFromJsonBytes(jsonFile []byte) (*Script, error) {
	script := &Script{}
	err := json.Unmarshal(jsonFile, script)
	return script, err
}
