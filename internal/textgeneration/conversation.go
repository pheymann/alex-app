package textgeneration

import "github.com/sashabaranov/go-openai"

type Conversation struct {
	Messages []Message `json:"messages" dynamodbav:"messages"`
}

type Message struct {
	Role string `json:"role" dynamodbav:"role"`
	Text string `json:"text" dynamodbav:"text"`
}

func (conversation *Conversation) AddPrompt(prompt string) {
	conversation.Messages = append(conversation.Messages, Message{
		Role: openai.ChatMessageRoleUser,
		Text: prompt,
	})
}
