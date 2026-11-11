// Package actions holds the operations shared by the CLI and the TUI. Each one
// strings together workspace (rules), render (bash) and store (files).
package actions

import (
	"github.com/sharon-xa/gomux/internal/render"
	"github.com/sharon-xa/gomux/internal/store"
	"github.com/sharon-xa/gomux/internal/workspace"
)

type Actions struct {
	store *store.Store
}

func New() *Actions {
	return &Actions{
		store: store.New(),
	}
}

// NewScript applies defaults, validates, then saves the spec and its
// rendered script. Validation errors are the workspace.Err* values.
func (a *Actions) NewScript(
	session workspace.Session,
	windows []workspace.Window,
	startingWindow int,
	panes []workspace.Pane,
) (*workspace.Script, error) {
	script := workspace.NewScript(session, windows, startingWindow, panes)
	script.ApplyDefaults()

	if err := script.Validate(); err != nil {
		return nil, err
	}

	if err := a.store.SaveSpec(script); err != nil {
		return nil, err
	}

	if err := a.store.SaveScript(script, []byte(render.Bash(script))); err != nil {
		return nil, err
	}

	return script, nil
}
