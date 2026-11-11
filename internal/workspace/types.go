package workspace

type Session struct {
	Name      string `json:"name"`      // Make sure to sanitaize this field, it will become a filename
	Directory string `json:"directory"` // the full path to where the session should start
	Type      string `json:"type"`      // Make sure to sanitaize this field, it will become a folder name
}

func NewSession(name, dir, sType string) *Session {
	return &Session{
		Name:      name,
		Directory: dir,
		Type:      sType,
	}
}

type Window struct {
	Name              string `json:"name"`
	Number            int    `json:"number"`
	Command           string `json:"command"`
	Layout            Layout `json:"layout"`
	SynchronizedPanes bool   `json:"synchronizedPanes"`
}

type Pane struct {
	ID        int       `json:"id"`
	Location  int       `json:"location"` // on which window will we create this pane
	Direction Direction `json:"direction"`
	Command   string    `json:"command"`
}

type Direction string

const (
	Horizontal Direction = "h"
	Vertical   Direction = "v"
)

type Layout string

const (
	LayoutTiled          Layout = "tiled"
	LayoutEvenHorizontal Layout = "even-horizontal"
	LayoutEvenVertical   Layout = "even-vertical"
	LayoutMainHorizontal Layout = "main-horizontal"
	LayoutMainVertical   Layout = "main-vertical"
)
