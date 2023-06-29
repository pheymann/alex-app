package talktome

import (
	"fmt"

	"talktome.com/internal/textgeneration"
)

type TalkToMe struct {
	textGen *textgeneration.TextGenerator
}

func NewTalkToMe(openAIToken string) TalkToMe {
	return TalkToMe{
		textGen: textgeneration.NewOpenAIGenerator(openAIToken),
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
		return nil, fmt.Errorf("failed to generate description:\n%w", err)
	}

	fmt.Printf("[DEBUG] Generate tasks for %s's \"%s\"\n", artistName, artName)

	tasks, err := talktome.textGen.GenerateTasks(artistName, artName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks:\n%w", err)
	}

	return &ArtPresentation{
		Description: description,
		Tasks:       tasks,
	}, nil
}
