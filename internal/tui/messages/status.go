package messages

type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusError
	StatusWarning
)

type StatusMsg struct {
	Message string
	Type    StatusType
}

func NewStatus(message string, statusType StatusType) StatusMsg {
	return StatusMsg{
		Message: message,
		Type:    statusType,
	}
}
