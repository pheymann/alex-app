package conversation

import (
	"time"

	"talktome.com/internal/idgenerator"
	"talktome.com/internal/textgeneration"
)

type Conversation struct {
	ID        string            `json:"id" yaml:"id" dynamodbav:"id"`
	Metadata  map[string]string `json:"metadata" yaml:"metadata" dynamodbav:"metadata"`
	Messages  []Message         `json:"messages" yaml:"messages" dynamodbav:"messages"`
	CreatedAt time.Time         `json:"createdAt" yaml:"createdAt" dynamodbav:"createdAt"`
	State     State             `json:"state" yaml:"state" dynamodbav:"state"`
}

type Message struct {
	Text                     string     `json:"text" yaml:"text" dynamodbav:"text"`
	Role                     string     `json:"role" yaml:"role" dynamodbav:"role"`
	SpeechClipUUID           string     `json:"speechClipUuid" yaml:"speechClipUuid" dynamodbav:"speech_clip_uuid"`
	SpeechClipURL            string     `json:"speechClipUrl" yaml:"speechClipUrl" dynamodbav:"speech_clip_url"`
	SpeechClipExpirationDate *time.Time `json:"speechClipExpirationDate" dynamodbav:"speech_clip_expiration_date"`
	SpeechClipIsExpired      bool       `json:"speechClipIsExpired" dynamodbav:"speech_clip_is_expired"`
}

type Role = string

const (
	RoleUser      Role = "user"
	RoleAssistent Role = "assistant"
	RoleSystem    Role = "system"
)

type State = string

const (
	StateGenerating State = "generating"
	StateReady      State = "ready"
)

func NewConversation(metadata map[string]string, idGen idgenerator.IDGenerator) Conversation {
	return Conversation{
		ID:        idGen.GenerateID(metadata),
		Metadata:  metadata,
		Messages:  []Message{},
		CreatedAt: time.Now(),
	}
}

func (conversation *Conversation) AddMessage(text string) {
	conversation.Messages = append(conversation.Messages, Message{
		Role: textgeneration.RoleUser,
		Text: text,
	})
}

func (conversation *Conversation) AddBasicMessage(message textgeneration.BasicMessage) {
	conversation.Messages = append(conversation.Messages, Message{
		Role: message.Role,
		Text: message.Text,
	})
}

func (conversation Conversation) ToBasicMessageHistory() []textgeneration.BasicMessage {
	basicMessages := []textgeneration.BasicMessage{}

	for _, message := range conversation.Messages {
		basicMessages = append(basicMessages, textgeneration.BasicMessage{
			Role: message.Role,
			Text: message.Text,
		})
	}

	return basicMessages
}
