package speechgeneration

import (
	"github.com/aws/aws-sdk-go/service/polly"
)

type SpeechGenerator struct {
	client       *polly.Polly
	engine       string
	outputFormat string
	englishVoice voice
	germanVoice  voice
}

type voice struct {
	male   string
	female string
}

func NewPollySpeechGenerator(pollyClient *polly.Polly) *SpeechGenerator {
	return &SpeechGenerator{
		client:       pollyClient,
		engine:       polly.EngineNeural,
		outputFormat: polly.OutputFormatMp3,
		englishVoice: voice{
			// the quality of Matthew is better than Ruth
			male:   polly.VoiceIdMatthew,
			female: polly.VoiceIdRuth,
		},
		germanVoice: voice{
			male: polly.VoiceIdDaniel,
			// quality isn't the best
			female: polly.VoiceIdVicki,
		},
	}
}
