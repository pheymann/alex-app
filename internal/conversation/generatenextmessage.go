package conversation

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func (generator *OpenAITextGenerationService) GenerateNextMessage(conversation *Conversation) error {
	openAIConversation := []openai.ChatCompletionMessage{}

	for _, message := range conversation.Messages {
		role := ""

		switch message.Role {
		case RoleAssistent:
			role = openai.ChatMessageRoleAssistant
		case RoleUser:
			role = openai.ChatMessageRoleUser
		case RoleSystem:
			role = openai.ChatMessageRoleAssistant
		default:
			return fmt.Errorf("unknown role: %s", message.Role)
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
		return fmt.Errorf("failed to generate text: %w", err)
	}

	conversation.Messages = append(conversation.Messages, Message{
		Role:        RoleAssistent,
		Text:        resp.Choices[0].Message.Content,
		CanHaveClip: true,
	})

	return nil
}
