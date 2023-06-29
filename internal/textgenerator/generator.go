package textgeneration

import openai "github.com/sashabaranov/go-openai"

type textGenerator struct {
	client *openai.Client
	model  string
}

func NewOpenAIGenerator(token string) *textGenerator {
	client := openai.NewClient(token)

	return &textGenerator{
		client: client,
		model:  openai.GPT3Dot5Turbo,
	}
}
