package textgeneration

import openai "github.com/sashabaranov/go-openai"

type OpenAITextGenerationService struct {
	client *openai.Client
	model  string
}

func NewOpenAIGenerator(token string) *OpenAITextGenerationService {
	client := openai.NewClient(token)

	return &OpenAITextGenerationService{
		client: client,
		model:  openai.GPT3Dot5Turbo,
	}
}
