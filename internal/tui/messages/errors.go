package messages

type ErrorMsg struct {
	Err error
}

func (e ErrorMsg) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func NewError(err error) ErrorMsg {
	return ErrorMsg{Err: err}
}
