package talktome

import (
	"fmt"

	"github.com/resemble-ai/resemble-go/v2/response"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
)

type TalkToMe struct {
	textGen   *textgeneration.TextGenerator
	speechGen *speechgeneration.SpeechGenerator
}

func NewTalkToMe(openAIToken string, resembleToken string, resembleProjectUUID string, resembleCallBackURL string) TalkToMe {
	return TalkToMe{
		textGen:   textgeneration.NewOpenAIGenerator(openAIToken),
		speechGen: speechgeneration.NewResembleGenerator(resembleToken, resembleProjectUUID, resembleCallBackURL),
	}
}

type ArtPresentation struct {
	Description string
	Tasks       []string
}

func (talktome TalkToMe) GenerateArtPresentation(artistName string, artName string) (*ArtPresentation, error) {
	fmt.Printf("[DEBUG] Generate description for %s's \"%s\"\n", artistName, artName)

	description, err := talktome.textGen.GenerateArtDescription(artistName, artName)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] Generate tasks for %s's \"%s\"\n", artistName, artName)

	tasks, err := talktome.textGen.GenerateTasks(artistName, artName)
	if err != nil {
		return nil, err
	}

	return &ArtPresentation{
		Description: description,
		Tasks:       tasks,
	}, nil
}

func (talktome TalkToMe) GenerateSpeechClip(text string) (response.Clip, error) {
	fmt.Printf("[DEBUG] Generate speech clip for %s\n", text)
	return talktome.speechGen.GenerateSpeechClip(text)
}

func (talktome TalkToMe) DownloadSpeechClip(clipURL string) ([]byte, error) {
	fmt.Printf("[DEBUG] Download speech clip %s\n", clipURL)
	return talktome.speechGen.DownloadSpeechClip(clipURL)
}
