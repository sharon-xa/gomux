package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
	"github.com/sharon-xa/gomux/internal/gomuxerr"
	"github.com/sharon-xa/gomux/internal/workspace"
)

type Actions struct {
	DB *database.DB
}

type IActions interface {
	NewScript(
		session workspace.Session,
		windows []workspace.Window,
		startingWindow int,
		panes []workspace.Pane,
	) (*workspace.Script, error)
}

func NewActions() IActions {
	return &Actions{
		DB: database.NewDB(),
	}
}

func (a *Actions) NewScript(
	session workspace.Session,
	windows []workspace.Window,
	startingWindow int,
	panes []workspace.Pane,
) (*workspace.Script, error) {
	if session.Name == "" {
		return nil, gomuxerr.ErrNoSessionNameFound
	}

	if session.Directory == "" {
		return nil, gomuxerr.ErrNoSessionDirectoryFound
	}

	if session.Type == "" {
		return nil, gomuxerr.ErrNoSessionTypeFound
	}

	if len(windows) == 0 {
		return nil, gomuxerr.ErrNoWindowFound
	}

	for _, window := range windows {
		if window.Layout == "" {
			window.Layout = workspace.LayoutEvenVertical
		}

		if window.Command == "" {
			return nil, gomuxerr.ErrNoWindowCommandFound
		}
	}

	if startingWindow <= 0 {
		startingWindow = 1
	}

	for _, pane := range panes {
		if pane.Direction == "" {
			pane.Direction = workspace.Vertical
		}
		if pane.Location == 0 {
			return nil, gomuxerr.ErrNoPaneLocationFound
		}
	}

	script := workspace.NewScript(session, windows, startingWindow, panes)
	err := a.DB.AddScript(script)
	if err != nil {
		return nil, err
	}

	return script, nil
}
