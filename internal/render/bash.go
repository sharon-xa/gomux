// Package render turns a workspace spec into a bash script that builds the
// tmux session. This is where roadmap milestones M1–M5 happen; the tests in
// this package define what "correct" means.
package render

import (
	"fmt"
	"strings"

	"github.com/sharon-xa/gomux/internal/workspace"
)

const (
	sectionSeparator       = "#" + "\n" + "--------" + "\n"
	sessionNameVariable    = "SESSION_NAME"
	sessionDirVariable     = "SESSION_DIR"
	scriptTemplateSections = 10
)

// Bash renders s as a bash script.
func Bash(s *workspace.Script) string {
	sections := make([]string, scriptTemplateSections)
	sections[0] = shebang()
	sections[1] = scriptVars(s.Session)
	sections[2] = checkDirectoryExistence()
	sections[3] = checkIfSessionExists()
	sections[4] = createSession()
	sections[5] = createWindows(s.Windows)
	sections[6] = createPanes(s.Panes)
	sections[7] = setStartingWindow(s.StartingWindow)

	return strings.Join(sections, sectionSeparator)
}

func shebang() string {
	return `#!/usr/bin/env bash`
}

func scriptVars(session workspace.Session) string {
	return fmt.Sprintf("%s=\"%s\"\n%s=\"%s\"",
		sessionNameVariable,
		session.Name,
		sessionDirVariable,
		session.Directory,
	)
}

func checkDirectoryExistence() string {
	return fmt.Sprintf(`if [ ! -d "$%s" ]; then
  echo "Directory $SESSION_DIR does not exist. Exiting..."
  exit 1
fi`, sessionDirVariable)
}

func checkIfSessionExists() string {
	return fmt.Sprintf(`if tmux has-session -t "$%s" 2>/dev/null; then
  tmux attach -t "$%s"
  exit 0
fi`, sessionNameVariable, sessionNameVariable)
}

func createSession() string {
	return fmt.Sprintf(
		`tmux new-session -d -s $%s -c "$%s"`,
		sessionNameVariable,
		sessionDirVariable,
	)
}

func createWindows(windows []workspace.Window) string {
	out := ""

	for i, window := range windows {
		out += fmt.Sprintf("# Window Number %d", i)
		out += createWindow(window)
		out += "\n"
	}

	return out
}

func createWindow(w workspace.Window) string {
	target := fmt.Sprintf(`$%s:%d`,
		sessionNameVariable,
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

func createPanes(panes []workspace.Pane) string {
	out := ""

	for i, pane := range panes {
		out += fmt.Sprintf("## Pane Number %d", i)
		out += createPane(pane)
		out += "\n"
	}

	return out
}

func createPane(p workspace.Pane) string {
	target := fmt.Sprintf(`$%s:$%d.%d`,
		sessionNameVariable,
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

func setStartingWindow(startingWindow int) string {
	return fmt.Sprintf(`tmux select-window -t $%s:%d`, sessionNameVariable, startingWindow)
}
