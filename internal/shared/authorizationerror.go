package shared

import "fmt"

type AuthorizationError struct {
	Cause   error
	Message string
}

func (e *AuthorizationError) Error() string {
	causeStr := ""
	if e.Cause != nil {
		causeStr = e.Cause.Error()
	}

	return fmt.Sprintf("%s:\n%s", e.Message, causeStr)
}
