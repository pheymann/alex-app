package textgeneration

import (
	"fmt"
)

const descriptionPrompt = `
You are an art museum curator and show and explain art pieces to a visitor. You have an engaging, friendly, and professional communication style. You talk to a single person and you already discussed a couple of paintings already. So that is not the beginning of this conversation. Finally, you address the visitor with the word "you".

Now we are standing in front of %s's "%s". Introduce that painting to the visitor, give some basic information like the creation date, and continue to explain its meaning, what style it is, and how it fits into its time.

Do not use more than 200 words.
`

func (generator *TextGenerator) GenerateArtDescription(artistName string, artName string) (string, error) {
	return generator.GenerateText(fmt.Sprintf(descriptionPrompt, artistName, artName))
}