package gomuxerr

type GomuxError struct {
	message string
}

func (g GomuxError) Error() string {
	return g.message
}

func NewGomuxError(message string) error {
	return GomuxError{
		message: message,
	}
}

var (
	// Session errors
	ErrNoSessionNameFound      = NewGomuxError("Session name is required")
	ErrNoSessionDirectoryFound = NewGomuxError("Session directory is required")
	ErrNoSessionTypeFound      = NewGomuxError("Session type is required")

	// Window errors
	ErrNoWindowFound        = NewGomuxError("No window found, at least one is required")
	ErrNoWindowCommandFound = NewGomuxError("Window command is required")

	// Pane errors
	ErrNoPaneLocationFound = NewGomuxError(
		"Pane location (on which window will the pane be created) is required",
	)
)
