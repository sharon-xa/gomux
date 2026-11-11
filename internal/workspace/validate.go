package workspace

import (
	"errors"
	"fmt"
)

// Validation errors. Compare with errors.Is; Validate may wrap them with the
// index of the offending window or pane.
var (
	ErrSessionNameRequired      = errors.New("session name is required")
	ErrSessionDirectoryRequired = errors.New("session directory is required")
	ErrSessionTypeRequired      = errors.New("session type is required")
	ErrNoWindows                = errors.New("at least one window is required")
	ErrWindowCommandRequired    = errors.New("window command is required")
	ErrPaneLocationRequired     = errors.New("pane location (the window to create the pane in) is required")
)

// ApplyDefaults fills in optional fields that were left empty.
func (s *Script) ApplyDefaults() {
	for i := range s.Windows {
		if s.Windows[i].Layout == "" {
			s.Windows[i].Layout = LayoutEvenVertical
		}
	}

	for i := range s.Panes {
		if s.Panes[i].Direction == "" {
			s.Panes[i].Direction = Vertical
		}
	}

	// TODO(M5): assumes 1-based windows, but tmux defaults to base-index 0 (review 5.2).
	if s.StartingWindow <= 0 {
		s.StartingWindow = 1
	}
}

// Validate returns the first problem that would stop the session being built.
func (s *Script) Validate() error {
	if s.Session.Name == "" {
		return ErrSessionNameRequired
	}

	if s.Session.Directory == "" {
		return ErrSessionDirectoryRequired
	}

	if s.Session.Type == "" {
		return ErrSessionTypeRequired
	}

	if len(s.Windows) == 0 {
		return ErrNoWindows
	}

	for i, w := range s.Windows {
		if w.Command == "" {
			return fmt.Errorf("window %d (%q): %w", i, w.Name, ErrWindowCommandRequired)
		}
	}

	for i, p := range s.Panes {
		// TODO(M5): this also rejects panes on window 0, the first window under
		// tmux's default base-index (review 5.2).
		if p.Location == 0 {
			return fmt.Errorf("pane %d: %w", i, ErrPaneLocationRequired)
		}
	}

	return nil
}
