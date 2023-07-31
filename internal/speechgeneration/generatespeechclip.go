package speechgeneration

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
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

	randomHash := rand.Intn(100)
	clipName := hex.EncodeToString([]byte(title + fmt.Sprint(randomHash)))
	outFile, err := os.Create(clipName + ".mp3")
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
