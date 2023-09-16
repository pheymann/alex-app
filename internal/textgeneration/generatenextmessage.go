package textgeneration

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	openai "github.com/sashabaranov/go-openai"
	"talktome.com/internal/shared"
)

func (generator *OpenAITextGenerationService) GenerateNextMessage(
	messageHistory []BasicMessage,
	logCtx zerolog.Context,
) (*BasicMessage, error) {
	shared.GetLogger(logCtx).Debug().Msgf("generate next message for history of %d messages", len(messageHistory))

	openAIConversation := []openai.ChatCompletionMessage{}

	for _, message := range messageHistory {
		role := ""

		switch message.Role {
		case RoleAssistent:
			role = openai.ChatMessageRoleAssistant
		case RoleUser:
			role = openai.ChatMessageRoleUser
		case RoleSystem:
			role = openai.ChatMessageRoleAssistant
		default:
			return nil, fmt.Errorf("unknown role: %s", message.Role)
		}

		openAIConversation = append(openAIConversation, openai.ChatCompletionMessage{
			Role:    role,
			Content: message.Text,
		})
	}

	resp, err := generator.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    generator.model,
			Messages: openAIConversation,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	return &BasicMessage{
		Role: RoleAssistent,
		Text: resp.Choices[0].Message.Content,
	}, nil
}
