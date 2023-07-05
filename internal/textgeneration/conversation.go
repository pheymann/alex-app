package textgeneration

import "github.com/sashabaranov/go-openai"

type Conversation struct {
	Messages []Message
}

type Message struct {
	Role string
	Text string
}

func (conversation *Conversation) AddPrompt(prompt string) {
	conversation.Messages = append(conversation.Messages, Message{
		Role: openai.ChatMessageRoleUser,
		Text: prompt,
	})
}
