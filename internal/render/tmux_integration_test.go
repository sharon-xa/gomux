package render_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sharon-xa/gomux/internal/render"
)

// TestTmuxIntegration runs each generated script against a real, isolated tmux
// server and checks the resulting session: windows, pane counts, working
// directories, the active window, and that every command ran in its own pane.
//
// Isolation: each case gets its own TMUX_TMPDIR (so its own server socket), its
// own HOME/XDG_CONFIG_HOME (so your ~/.config/tmux/tmux.conf is NOT loaded),
// SHELL=/bin/sh for predictable pane shells, and TMUX unset.
//
// Contract for the generator: when GOMUX_NO_ATTACH=1 is set, the script must
// not attach or switch-client at the end (tests have no terminal).
//
// Skipped with -short (`just test-unit`) or when tmux isn't installed.
func TestTmuxIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tmux integration test in -short mode")
	}
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		t.Skip("tmux not found in PATH")
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			t.Parallel()
			runFixture(t, tmux, f)
		})
	}
}

func runFixture(t *testing.T, tmuxBin string, f fixture) {
	// Short root path: tmux socket paths must fit in ~100 bytes.
	root, err := os.MkdirTemp("", "gmx")
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(root, "home")
	sockDir := filepath.Join(root, "sock")
	projDir := filepath.Join(root, f.dirName)
	for _, d := range []string{home, sockDir, projDir} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}

	if f.tmuxConf != "" {
		if err := os.WriteFile(filepath.Join(home, ".tmux.conf"), []byte(f.tmuxConf), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	env := isolatedEnv(home, sockDir)
	tmux := func(args ...string) (string, error) {
		cmd := exec.Command(tmuxBin, args...)
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		return strings.TrimRight(string(out), "\n"), err
	}

	script := render.Bash(f.build(projDir))
	scriptPath := filepath.Join(root, "script.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = tmux("kill-server")
		if t.Failed() {
			t.Logf("--- generated script ---\n%s", numbered(script))
		}
		_ = os.RemoveAll(root)
	})

	// Run the script directly (not via `bash script.sh`) so the shebang is tested too.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var out []byte
	for attempt := 0; ; attempt++ {
		run := exec.CommandContext(ctx, scriptPath)
		run.Env = env
		run.Dir = root
		out, err = run.CombinedOutput()
		// ETXTBSY: a parallel test's fork briefly inherited our write fd for the
		// script (golang/go#22315). Harmless; just retry.
		if !errors.Is(err, syscall.ETXTBSY) || attempt == 10 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("script timed out (is it trying to attach despite GOMUX_NO_ATTACH=1?)\n%s", out)
	}
	if err != nil {
		t.Fatalf("script failed: %v\n%s", err, out)
	}

	w := f.want
	target := "=" + w.session

	if out, err := tmux("has-session", "-t", target); err != nil {
		t.Fatalf("session %q was not created: %s", w.session, out)
	}

	// Windows, in order.
	gotWindows, err := tmux("list-windows", "-t", target, "-F", "#{window_name}")
	if err != nil {
		t.Fatalf("list-windows: %v: %s", err, gotWindows)
	}
	var wantNames []string
	for _, ww := range w.windows {
		wantNames = append(wantNames, ww.name)
	}
	if got := strings.Split(gotWindows, "\n"); !slices.Equal(got, wantNames) {
		t.Errorf("windows = %q, want %q", got, wantNames)
	}

	// Pane counts and working directories.
	panes, err := tmux("list-panes", "-s", "-t", target, "-F", "#{window_name}\t#{pane_id}\t#{pane_current_path}")
	if err != nil {
		t.Fatalf("list-panes: %v: %s", err, panes)
	}
	wantDir, _ := filepath.EvalSymlinks(projDir)
	paneCount := map[string]int{}
	var paneIDs []string
	for line := range strings.SplitSeq(panes, "\n") {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}
		paneCount[parts[0]]++
		paneIDs = append(paneIDs, parts[1])
		if gotDir, _ := filepath.EvalSymlinks(parts[2]); gotDir != wantDir {
			t.Errorf("pane %s in window %q started in %q, want %q", parts[1], parts[0], parts[2], projDir)
		}
	}
	for _, ww := range w.windows {
		if paneCount[ww.name] != ww.panes {
			t.Errorf("window %q has %d panes, want %d", ww.name, paneCount[ww.name], ww.panes)
		}
	}

	// Active window. (display-message prints nothing without an attached
	// client, so read the window_active flag instead.)
	if list, err := tmux("list-windows", "-t", target, "-F", "#{window_active}\t#{window_name}"); err != nil {
		t.Errorf("list-windows: %v: %s", err, list)
	} else {
		var active string
		for line := range strings.SplitSeq(list, "\n") {
			if name, ok := strings.CutPrefix(line, "1\t"); ok {
				active = name
			}
		}
		if active != w.activeWindow {
			t.Errorf("active window = %q, want %q", active, w.activeWindow)
		}
	}

	// Commands ran, in the session directory, with their text intact.
	for name, want := range w.files {
		path := filepath.Join(projDir, name)
		got, ok := waitForFile(path, want, 10*time.Second)
		switch {
		case !ok && got == "":
			t.Errorf("%s was never written (command didn't run, or ran in the wrong directory)", name)
		case !ok:
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	// Each command was typed into its own pane.
	if len(w.ownPane) > 0 {
		owners := map[string][]string{}
		for _, id := range paneIDs {
			content, _ := tmux("capture-pane", "-p", "-t", id)
			for _, marker := range w.ownPane {
				if strings.Contains(content, marker) {
					owners[marker] = append(owners[marker], id)
				}
			}
		}
		seen := map[string]string{}
		for _, marker := range w.ownPane {
			ids := owners[marker]
			if len(ids) != 1 {
				t.Errorf("command for %q appears in %d panes %v, want exactly 1", marker, len(ids), ids)
				continue
			}
			if other, dup := seen[ids[0]]; dup {
				t.Errorf("commands for %q and %q were both sent to pane %s", other, marker, ids[0])
			}
			seen[ids[0]] = marker
		}
	}
}

// isolatedEnv returns the current environment with everything that could
// leak the user's tmux server, config or shell setup replaced.
func isolatedEnv(home, sockDir string) []string {
	drop := []string{"TMUX", "TMUX_PANE", "TMUX_TMPDIR", "HOME", "XDG_CONFIG_HOME", "SHELL", "ENV", "BASH_ENV", "GOMUX_NO_ATTACH"}
	var env []string
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		if !slices.Contains(drop, key) {
			env = append(env, kv)
		}
	}
	return append(env,
		"HOME="+home,
		"XDG_CONFIG_HOME="+filepath.Join(home, ".config"),
		"TMUX_TMPDIR="+sockDir,
		"SHELL=/bin/sh",
		"GOMUX_NO_ATTACH=1",
	)
}

// waitForFile polls until path has the wanted content or the timeout passes.
// Pane commands run asynchronously after send-keys, so this can't be a single read.
func waitForFile(path, want string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)
	var got string
	for {
		if b, err := os.ReadFile(path); err == nil {
			got = string(b)
			if got == want {
				return got, true
			}
		}
		if time.Now().After(deadline) {
			return got, false
		}
		time.Sleep(50 * time.Millisecond)
	}
}
