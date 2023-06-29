package textgeneration

import openai "github.com/sashabaranov/go-openai"

type TextGenerator struct {
	client *openai.Client
	model  string
}

func NewOpenAIGenerator(token string) *TextGenerator {
	client := openai.NewClient(token)

	return &TextGenerator{
		client: client,
		model:  openai.GPT3Dot5Turbo,
	}
}
