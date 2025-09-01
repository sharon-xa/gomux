package models

import (
	"fmt"
	"strings"
)

type Session struct {
	Name      string `json:"name"`      // Make sure to sanitaize this field, it will become a filename
	Directory string `json:"directory"` // the full path to where the session should start
	Type      string `json:"type"`      // Make sure to sanitaize this field, it will become a folder name
}

func NewSession(name, dir, typ string) *Session {
	return &Session{
		Name:      name,
		Directory: dir,
		Type:      typ,
	}
}

func (s *Session) createSession() string {
	return fmt.Sprintf(
		`tmux new-session -d -s $%s -c "$%s"`,
		SESSION_NAME_VARIABLE,
		SESSION_DIR_VARIABLE,
	)
}

type Window struct {
	Name              string `json:"name"`
	Number            int    `json:"number"`
	Command           string `json:"command"`
	Layout            Layout `json:"layout"`
	SynchronizedPanes bool   `json:"synchronizedPanes"`
}

func (w *Window) createWindow() string {
	target := fmt.Sprintf(`$%s:%d`,
		SESSION_NAME_VARIABLE,
		w.Number,
	)

	var builder strings.Builder

	baseCommand := fmt.Sprintf(
		`tmux rename-window -t %s "%s"
tmux send-keys -t %s "%s" C-m`,
		target,
		w.Name,
		target,
		w.Command,
	)
	builder.WriteString(baseCommand)

	if w.Layout != "" {
		builder.WriteString("\n")
		windowLayoutCommand := fmt.Sprintf(
			`tmux select-layout -t %s %s`,
			target,
			w.Layout,
		)
		builder.WriteString(windowLayoutCommand)
	}

	if w.SynchronizedPanes {
		builder.WriteString("\n")
		synchronizedPanesCommand := fmt.Sprintf(
			`tmux set-window-option -t %s synchronize-panes on`,
			target,
		)
		builder.WriteString(synchronizedPanesCommand)
	}

	return builder.String()
}

type Pane struct {
	ID        int       `json:"id"`
	Location  int       `json:"location"` // on which window will we create this pane
	Direction Direction `json:"direction"`
	Command   string    `json:"command"`
}

func (p *Pane) createPane() string {
	target := fmt.Sprintf(`$%s:$%d.%d`,
		SESSION_NAME_VARIABLE,
		p.Location,
		p.ID,
	)

	return fmt.Sprintf(`tmux split-window -%s -t %s
tmux send-keys -t %s "%s" C-m`,
		p.Direction,
		target,
		target,
		p.Command,
	)
}
