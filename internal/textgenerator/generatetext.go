package textgeneration

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

func (generator *textGenerator) GenerateText(prompt string) (string, error) {
	resp, err := generator.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: generator.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
