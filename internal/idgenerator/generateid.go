package idgenerator

import (
	"encoding/base64"
	"time"
)

func (generator *randomIDGenerator) GenerateID(metadata map[string]string) string {
	metadataStr := time.Now().Format(time.RFC3339Nano)
	for _, value := range metadata {
		metadataStr += "::" + value
	}

	return base64.StdEncoding.EncodeToString([]byte(metadataStr))
}
