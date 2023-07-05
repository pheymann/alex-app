package textgeneration

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func (generator *TextGenerator) ContinueConversation(conversation *Conversation) error {
	openAIConversation := []openai.ChatCompletionMessage{}

	for _, message := range conversation.Messages {
		openAIConversation = append(openAIConversation, openai.ChatCompletionMessage{
			Role:    message.Role,
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
		Role: openai.ChatMessageRoleAssistant,
		Text: resp.Choices[0].Message.Content,
	})

	return nil
}
