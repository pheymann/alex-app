package textgeneration

import openai "github.com/sashabaranov/go-openai"

type TextGenerator interface {
	ContinueConversation(*Conversation) error
}

type OpenAITextGenerator struct {
	client *openai.Client
	model  string
}

func NewOpenAIGenerator(token string) *OpenAITextGenerator {
	client := openai.NewClient(token)

	return &OpenAITextGenerator{
		client: client,
		model:  openai.GPT3Dot5Turbo,
	}
}
