package shared

import "fmt"

type InternalError struct {
	Cause   error
	Message string
}

func (e *InternalError) Error() string {
	causeStr := ""
	if e.Cause != nil {
		causeStr = e.Cause.Error()
	}

	return fmt.Sprintf("%s:\n%s", e.Message, causeStr)
}
