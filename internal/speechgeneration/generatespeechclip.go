package speechgeneration

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"talktome.com/internal/shared"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
)

func (generator *AWSPollySpeechGenerator) GenerateSpeechClip(
	title string,
	text string,
	language shared.Language,
	logCtx zerolog.Context,
) (*os.File, error) {
	shared.GetLogger(logCtx).Debug().Msg("generate speech clip")

	var voiceId string
	switch language {
	case shared.LanguageGerman:
		voiceId = generator.germanVoice.male
	case shared.LanguageEnglish:
		voiceId = generator.englishVoice.male
	default:
		return nil, fmt.Errorf("unknown language: %s", language)
	}

	resp, err := generator.client.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
		Engine:       &generator.engine,
		OutputFormat: &generator.outputFormat,
		Text:         aws.String(text),
		VoiceId:      &voiceId,
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
