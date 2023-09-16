package entitystore

import "fmt"

type EntityStoreError struct {
	Cause   error
	Message string
}

func (e *EntityStoreError) Error() string {
	return fmt.Sprintf("%s:\n%s", e.Message, e.Cause.Error())
}

func NewEntityStoreError(message string, cause error) *EntityStoreError {
	return &EntityStoreError{
		Cause:   cause,
		Message: message,
	}
}
