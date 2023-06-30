package speechgeneration

import (
	"net/http"

	resemble "github.com/resemble-ai/resemble-go/v2"
)

type SpeechGenerator struct {
	client      *resemble.Client
	httpClient  *http.Client
	voiceUUID   string
	projectUUID string
	callbackURL string
}

func NewResembleGenerator(token string, projectUUID string, callbackURL string) *SpeechGenerator {
	client := resemble.NewClient(token)

	return &SpeechGenerator{
		client:      client,
		httpClient:  &http.Client{},
		projectUUID: projectUUID,
		callbackURL: callbackURL,
	}
}
