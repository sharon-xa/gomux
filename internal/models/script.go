package models

type Script struct {
	ID             int
	Session        Session
	Window         []Window
	StartingWindow int // when the tmux start on which window should we start
	Pane           []Pane
}

func (s *Script) Stringify() string {
	return ""
}
