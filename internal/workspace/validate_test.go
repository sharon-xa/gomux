package workspace

import (
	"errors"
	"testing"
)

func validScript() *Script {
	return NewScript(
		Session{Name: "dev", Directory: "/tmp", Type: "work"},
		[]Window{{Name: "editor", Number: 1, Command: "nvim ."}},
		1,
		nil,
	)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(s *Script)
		wantErr error
	}{
		{
			name:    "valid",
			modify:  func(s *Script) {},
			wantErr: nil,
		},
		{
			name:    "missing session name",
			modify:  func(s *Script) { s.Session.Name = "" },
			wantErr: ErrSessionNameRequired,
		},
		{
			name:    "missing session directory",
			modify:  func(s *Script) { s.Session.Directory = "" },
			wantErr: ErrSessionDirectoryRequired,
		},
		{
			name:    "missing session type",
			modify:  func(s *Script) { s.Session.Type = "" },
			wantErr: ErrSessionTypeRequired,
		},
		{
			name:    "no windows",
			modify:  func(s *Script) { s.Windows = nil },
			wantErr: ErrNoWindows,
		},
		{
			name: "window without command",
			modify: func(s *Script) {
				s.Windows = append(s.Windows, Window{Name: "server", Number: 2})
			},
			wantErr: ErrWindowCommandRequired,
		},
		{
			name: "pane without location",
			modify: func(s *Script) {
				s.Panes = []Pane{{Command: "htop"}}
			},
			wantErr: ErrPaneLocationRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validScript()
			tt.modify(s)

			err := s.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNamesTheOffendingWindow(t *testing.T) {
	s := validScript()
	s.Windows = append(s.Windows, Window{Name: "server", Number: 2})

	err := s.Validate()
	if err == nil || err.Error() != `window 1 ("server"): window command is required` {
		t.Errorf("Validate() = %v", err)
	}
}

func TestApplyDefaults(t *testing.T) {
	s := NewScript(
		Session{Name: "dev", Directory: "/tmp", Type: "work"},
		[]Window{
			{Name: "editor", Command: "nvim ."},
			{Name: "server", Command: "npm run dev", Layout: LayoutTiled},
		},
		0,
		[]Pane{
			{Location: 1, Command: "htop"},
			{Location: 1, Command: "top", Direction: Horizontal},
		},
	)

	s.ApplyDefaults()

	if got := s.Windows[0].Layout; got != LayoutEvenVertical {
		t.Errorf("empty layout = %q, want default %q", got, LayoutEvenVertical)
	}
	if got := s.Windows[1].Layout; got != LayoutTiled {
		t.Errorf("explicit layout = %q, want it kept as %q", got, LayoutTiled)
	}
	if got := s.Panes[0].Direction; got != Vertical {
		t.Errorf("empty direction = %q, want default %q", got, Vertical)
	}
	if got := s.Panes[1].Direction; got != Horizontal {
		t.Errorf("explicit direction = %q, want it kept as %q", got, Horizontal)
	}
	if s.StartingWindow != 1 {
		t.Errorf("StartingWindow = %d, want default 1", s.StartingWindow)
	}
}
