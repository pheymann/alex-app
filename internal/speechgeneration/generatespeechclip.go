package speechgeneration

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
)

func (generator *AWSPollySpeechGenerator) GenerateSpeechClip(title string, text string) (*os.File, error) {
	log.Debug().Msg("synthesize clip")
	resp, err := generator.client.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
		Engine:       &generator.engine,
		OutputFormat: &generator.outputFormat,
		Text:         aws.String(text),
		// TODO: make this a selection
		VoiceId: &generator.englishVoice.male,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to send Polly request: %w", err)
	}

	outFile, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary mp3 file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.AudioStream)
	if err != nil {
		return nil, fmt.Errorf("failed to download speech clip: %w", err)
	}

	return os.Open(outFile.Name())
}
