package render_test

import "github.com/sharon-xa/gomux/internal/workspace"

// fixture is one end-to-end scenario for the script generator.
//
// build constructs the spec using the *current* model shape. When the model
// changes (e.g. panes nested under windows), only the build funcs need to
// change: the expectations below describe tmux behaviour, not model fields.
type fixture struct {
	name string

	// dirName is the session directory created inside the test's temp dir.
	dirName string

	build func(dir string) *workspace.Script

	// tmuxConf, if set, is written to the isolated ~/.tmux.conf before the
	// script runs, to check the script doesn't depend on the user's config.
	tmuxConf string

	want want
}

type want struct {
	session string

	// windows in order, by name, with their expected pane count.
	windows []windowWant

	// activeWindow is the window that should be selected after the script runs.
	activeWindow string

	// files are written by the window/pane commands, relative to the session
	// directory. They prove the commands actually ran, in the right directory,
	// with their text passed through unchanged.
	files map[string]string

	// ownPane: each marker must appear in exactly one pane, and every marker in
	// a different pane. Proves each command was sent to its own pane.
	ownPane []string
}

type windowWant struct {
	name  string
	panes int
}

// goldenDir is the fixed session directory used when rendering golden files,
// so the output doesn't depend on where the tests run.
const goldenDir = "/home/user/projects/gomux-golden"

var fixtures = []fixture{
	{
		// Milestone 1: the smallest script that has to work.
		name:    "single-window",
		dirName: "project",
		build: func(dir string) *workspace.Script {
			return workspace.NewScript(
				workspace.Session{Name: "gomux-single", Directory: dir, Type: "test"},
				[]workspace.Window{
					{Name: "editor", Number: 0, Command: "echo ready > single.txt"},
				},
				0,
				nil,
			)
		},
		want: want{
			session:      "gomux-single",
			windows:      []windowWant{{"editor", 1}},
			activeWindow: "editor",
			files:        map[string]string{"single.txt": "ready\n"},
		},
	},
	{
		// Milestone 2: several windows, starting on one that isn't the first.
		name:    "multi-window",
		dirName: "project",
		build: func(dir string) *workspace.Script {
			return workspace.NewScript(
				workspace.Session{Name: "gomux-multi", Directory: dir, Type: "test"},
				[]workspace.Window{
					{Name: "editor", Number: 0, Command: "echo editor > editor.txt"},
					{Name: "server", Number: 1, Command: "echo server > server.txt"},
					{Name: "logs", Number: 2, Command: "echo logs > logs.txt"},
				},
				1,
				nil,
			)
		},
		want: want{
			session:      "gomux-multi",
			windows:      []windowWant{{"editor", 1}, {"server", 1}, {"logs", 1}},
			activeWindow: "server",
			files: map[string]string{
				"editor.txt": "editor\n",
				"server.txt": "server\n",
				"logs.txt":   "logs\n",
			},
			ownPane: []string{"editor.txt", "server.txt", "logs.txt"},
		},
	},
	{
		// Milestone 3: one window split into three panes, each running its own command.
		name:    "panes",
		dirName: "project",
		build: func(dir string) *workspace.Script {
			return workspace.NewScript(
				workspace.Session{Name: "gomux-panes", Directory: dir, Type: "test"},
				[]workspace.Window{
					{Name: "dev", Number: 0, Command: "echo main > main.txt", Layout: workspace.LayoutTiled},
				},
				0,
				// Current model: ID = index of the pane being split, Location = window number.
				[]workspace.Pane{
					{ID: 0, Location: 0, Direction: workspace.Horizontal, Command: "echo right > right.txt"},
					{ID: 1, Location: 0, Direction: workspace.Vertical, Command: "echo bottom > bottom.txt"},
				},
			)
		},
		want: want{
			session:      "gomux-panes",
			windows:      []windowWant{{"dev", 3}},
			activeWindow: "dev",
			files: map[string]string{
				"main.txt":   "main\n",
				"right.txt":  "right\n",
				"bottom.txt": "bottom\n",
			},
			ownPane: []string{"main.txt", "right.txt", "bottom.txt"},
		},
	},
	{
		// Milestone 4: user text must reach the pane shell unchanged:
		// double quotes, single quotes, $ and backslashes, plus a directory with a space.
		name:    "quoting",
		dirName: "my project",
		build: func(dir string) *workspace.Script {
			return workspace.NewScript(
				workspace.Session{Name: "gomux-quoting", Directory: dir, Type: "test"},
				[]workspace.Window{
					{
						Name:    "shell",
						Number:  0,
						Command: `printf '%s\n' "double \"quoted\"" 'single $HOME' 'back\slash' > quoted.txt`,
					},
				},
				0,
				nil,
			)
		},
		want: want{
			session:      "gomux-quoting",
			windows:      []windowWant{{"shell", 1}},
			activeWindow: "shell",
			files: map[string]string{
				"quoted.txt": "double \"quoted\"\nsingle $HOME\nback\\slash\n",
			},
		},
	},
	{
		// Milestone 5: the same spec must work under a common user config
		// (this is the setting in the author's own tmux.conf).
		name:    "user-config-base-index-1",
		dirName: "project",
		tmuxConf: "set -g base-index 1\n" +
			"set -g pane-base-index 1\n" +
			"set-window-option -g pane-base-index 1\n",
		build: func(dir string) *workspace.Script {
			return workspace.NewScript(
				workspace.Session{Name: "gomux-baseindex", Directory: dir, Type: "test"},
				[]workspace.Window{
					{Name: "editor", Number: 0, Command: "echo editor > editor.txt"},
					{Name: "server", Number: 1, Command: "echo server > server.txt"},
				},
				1,
				[]workspace.Pane{
					{ID: 0, Location: 0, Direction: workspace.Horizontal, Command: "echo split > split.txt"},
				},
			)
		},
		want: want{
			session:      "gomux-baseindex",
			windows:      []windowWant{{"editor", 2}, {"server", 1}},
			activeWindow: "server",
			files: map[string]string{
				"editor.txt": "editor\n",
				"server.txt": "server\n",
				"split.txt":  "split\n",
			},
			ownPane: []string{"editor.txt", "split.txt", "server.txt"},
		},
	},
}
