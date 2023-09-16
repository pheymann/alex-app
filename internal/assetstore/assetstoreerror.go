package assetstore

import "fmt"

type AsssetStoreError struct {
	Cause   error
	Message string
}

func (e *AsssetStoreError) Error() string {
	return fmt.Sprintf("%s:\n%s", e.Message, e.Cause.Error())
}
