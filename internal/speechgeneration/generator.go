package speechgeneration

import (
	"os"

	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/rs/zerolog"
)

type SpeechGenerator interface {
	GenerateSpeechClip(title string, text string, logCtx zerolog.Context) (*os.File, error)
}

type AWSPollySpeechGenerator struct {
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

func NewAWSPollySpeechGenerator(pollyClient *polly.Polly) *AWSPollySpeechGenerator {
	return &AWSPollySpeechGenerator{
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
