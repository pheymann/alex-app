package speechgeneration

import (
	"fmt"
	"io"

	"github.com/resemble-ai/resemble-go/v2/request"
	"github.com/resemble-ai/resemble-go/v2/response"
)

func (generator *SpeechGenerator) GenerateSpeechClip(text string) (response.Clip, error) {
	clip, err := generator.client.Clip.CreateAsync(generator.projectUUID, generator.callbackURL, request.Payload{
		"voice_uuid": generator.voiceUUID,
		"body":       text,
		"raw":        false,
	})

	if err != nil {
		return response.Clip{}, fmt.Errorf("failed to generate a clip: %w", err)
	}

	if !clip.Success {
		return response.Clip{}, fmt.Errorf("resemble.ai clip generation unsuccessful: %s", clip.Item.Body)
	}

	return clip, nil
}

func (generator *SpeechGenerator) DownloadSpeechClip(clipURL string) ([]byte, error) {
	response, err := generator.httpClient.Get(clipURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download clip %s: %w", clipURL, err)
	}

	defer response.Body.Close()

	clipBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read clip %s: %w", clipURL, err)
	}

	return clipBytes, nil
}
