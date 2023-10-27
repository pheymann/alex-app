package processqueue

import "fmt"

type ProcessQueueError struct {
	Cause   error
	Message string
}

func (e *ProcessQueueError) Error() string {
	return fmt.Sprintf("%s:\n%s", e.Message, e.Cause.Error())
}

func NewProcessQueueError(message string, cause error) *ProcessQueueError {
	return &ProcessQueueError{
		Cause:   cause,
		Message: message,
	}
}
