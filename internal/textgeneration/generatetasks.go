package textgeneration

import (
	"fmt"
	"strings"
)

const taskPrompt = `
You are an art museum curator and show and explain art pieces to a visitor. You have an engaging, friendly, and professional communication style. You talk to a single person and you already discussed a couple of paintings already. So that is not the beginning of this conversation. Finally, you address the visitor with the word "you".

You already introduced %s's "%s" and now you point out interesting facets of the painting the viewer should take note of or look at. Please don't provide more than 3 and split each hint into its own paragraph. Also don't ask if there are any more questions.
`

func (generator *TextGenerator) GenerateTasks(artistName string, artName string) ([]string, error) {
	taskString, err := generator.GenerateText(fmt.Sprintf(taskPrompt, artistName, artName))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks for %s's \"%s\": %w", artistName, artName, err)
	}

	return strings.Split(taskString, "\n\n"), nil
}
