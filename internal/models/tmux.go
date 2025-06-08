package models

type Session struct {
	ID               int
	SessionName      string
	SessionDirectory string
	SessionType      string
}

type Window struct {
	WindowName    string
	WindowNumber  int
	WindowCommand string
}

type Pane struct {
	PaneLocation  int // on which window will we create this pane
	PaneDirection string
	PaneCommand   string
	PanesLayout   string // for example: tiled
	Synchronized  bool
}
