package actions

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sharon-xa/gomux/internal/workspace"
)

// newTestActions points the data dir at a temp HOME so tests never touch
// your real ~/.local/share/gomux.
func newTestActions(t *testing.T) (*Actions, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	return New(), home
}

func validSession() workspace.Session {
	return workspace.Session{Name: "dev", Directory: "/tmp", Type: "work"}
}

func validWindows() []workspace.Window {
	return []workspace.Window{{Name: "editor", Number: 0, Command: "nvim ."}}
}

// The rules themselves are tested in workspace; this checks NewScript
// enforces them and writes nothing when they fail.
func TestNewScriptRejectsInvalidAndSavesNothing(t *testing.T) {
	a, home := newTestActions(t)

	script, err := a.NewScript(workspace.Session{Directory: "/tmp", Type: "work"}, validWindows(), 0, nil)
	if !errors.Is(err, workspace.ErrSessionNameRequired) {
		t.Fatalf("err = %v, want %v", err, workspace.ErrSessionNameRequired)
	}
	if script != nil {
		t.Errorf("script = %+v, want nil on error", script)
	}

	if _, err := os.Stat(filepath.Join(home, ".local")); !os.IsNotExist(err) {
		t.Errorf("files were written for an invalid script (stat err = %v)", err)
	}
}

// Review 5.1: defaults used to be set on loop copies and lost.
func TestNewScriptAppliesDefaults(t *testing.T) {
	a, _ := newTestActions(t)

	script, err := a.NewScript(
		validSession(),
		[]workspace.Window{
			{Name: "editor", Number: 0, Command: "nvim ."},
			{Name: "server", Number: 1, Command: "npm run dev"},
		},
		0,
		[]workspace.Pane{{ID: 0, Location: 1, Command: "htop"}},
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, w := range script.Windows {
		if w.Layout != workspace.LayoutEvenVertical {
			t.Errorf("window %q layout = %q, want default %q", w.Name, w.Layout, workspace.LayoutEvenVertical)
		}
	}
	for i, p := range script.Panes {
		if p.Direction != workspace.Vertical {
			t.Errorf("pane %d direction = %q, want default %q", i, p.Direction, workspace.Vertical)
		}
	}
}

// Review 4.1: the spec should be saved as <name>.tmux.json and the script as
// <name>.tmux.sh, under a folder named after the session type.
func TestNewScriptPersistsSpecAndScript(t *testing.T) {
	a, home := newTestActions(t)

	if _, err := a.NewScript(validSession(), validWindows(), 0, nil); err != nil {
		t.Fatal(err)
	}

	base := filepath.Join(home, ".local", "share", "gomux")
	spec := filepath.Join(base, "dataState", "work", "dev.tmux.json")
	script := filepath.Join(base, "data", "work", "dev.tmux.sh")

	if _, err := os.Stat(spec); err != nil {
		t.Errorf("spec not saved at %s: %v", spec, err)
	}

	info, err := os.Stat(script)
	if err != nil {
		t.Fatalf("script not saved at %s: %v", script, err)
	}
	if info.Mode().Perm()&0o100 == 0 {
		t.Errorf("script %s is not executable (mode %v)", script, info.Mode().Perm())
	}
}
