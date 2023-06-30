package textgeneration

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func (generator *TextGenerator) GenerateText(prompt string) (string, error) {
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
		return "", fmt.Errorf("failed to generate text: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
