package actions

import (
	"errors"

	"github.com/sharon-xa/gomux/internal/models"
)

func (a *Actions) CreateScriptAllInOne(script *models.Script) error {
	if len(script.Windows) == 0 {
		return errors.New("there have to be at least one window")
	}

	err := a.DB.AddScript(script)
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) CreateScript() error {
	if len(a.CurrentScript.Windows) == 0 {
		return errors.New("there have to be at least one window")
	}

	err := a.DB.AddScript(a.CurrentScript)
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) AddSession(s *models.Session) error {
	if s.Name == "" {
		return errors.New("Session name is required")
	}

	if s.Directory == "" {
		return errors.New("Session directory is required")
	}

	if s.Type == "" {
		return errors.New("Session type is required")
	}

	a.CurrentScript.Session = *s
	return nil
}

func (a *Actions) AddWindow(w *models.Window) error {
	if w.Layout == "" {
		w.Layout = models.LayoutEvenVertical
	}

	if w.Command == "" {
		return errors.New("Window command is required")
	}

	a.CurrentScript.Windows = append(a.CurrentScript.Windows, *w)
	return nil
}

func (a *Actions) AddPane(p *models.Pane) error {
	if p.Direction == "" {
		p.Direction = models.Vertical
	}

	a.CurrentScript.Panes = append(a.CurrentScript.Panes, *p)
	return nil
}
