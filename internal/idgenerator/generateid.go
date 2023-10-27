package idgenerator

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

func (generator *randomIDGenerator) GenerateID(metadata map[string]string) string {
	metadataStr := fmt.Sprint(rand.Intn(999_999_999))
	for _, value := range metadata {
		metadataStr += "::" + value
	}

	return base64.StdEncoding.EncodeToString([]byte(metadataStr))
}
