package shared

import "fmt"

type UserInputError struct {
	Message string
}

func (e *UserInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Message)
}
