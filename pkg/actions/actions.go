package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
	"github.com/sharon-xa/gomux/internal/gomuxerr"
	"github.com/sharon-xa/gomux/internal/models"
)

type Actions struct {
	DB *database.DB
}

type IActions interface {
	NewScript(
		session models.Session,
		windows []models.Window,
		startingWindow int,
		panes []models.Pane,
	) (*models.Script, error)
}

func NewActions() IActions {
	return &Actions{
		DB: database.NewDB(),
	}
}

func (a *Actions) NewScript(
	session models.Session,
	windows []models.Window,
	startingWindow int,
	panes []models.Pane,
) (*models.Script, error) {
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
			window.Layout = models.LayoutEvenVertical
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
			pane.Direction = models.Vertical
		}
		if pane.Location == 0 {
			return nil, gomuxerr.ErrNoPaneLocationFound
		}
	}

	script := models.NewScript(session, windows, startingWindow, panes)
	err := a.DB.AddScript(script)
	if err != nil {
		return nil, err
	}

	return script, nil
}
