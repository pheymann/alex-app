package conversation

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"talktome.com/internal/textgeneration"
)

type Conversation struct {
	ID       string            `json:"id" yaml:"id" dynamodbav:"id"`
	Metadata map[string]string `json:"metadata" yaml:"metadata" dynamodbav:"metadata"`
	Messages []Message         `json:"messages" yaml:"messages" dynamodbav:"messages"`
}

type Message struct {
	Text                     string     `json:"text" yaml:"text" dynamodbav:"text"`
	Role                     string     `json:"role" yaml:"role" dynamodbav:"role"`
	CanHaveClip              bool       `json:"canHaveClip" dynamodbav:"can_have_clip"`
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

func NewConversation(metadata map[string]string) Conversation {
	return Conversation{
		ID:       GenerateStableID(metadata),
		Metadata: metadata,
		Messages: []Message{},
	}
}

func GenerateStableID(metadata map[string]string) string {
	idSeed := fmt.Sprint(rand.Intn(1000))
	for _, value := range metadata {
		idSeed += "::" + value
	}

	return base64.StdEncoding.EncodeToString([]byte(idSeed))
}

func (conversation *Conversation) AddMessage(text string) {
	conversation.Messages = append(conversation.Messages, Message{
		Role:        textgeneration.RoleUser,
		Text:        text,
		CanHaveClip: false,
	})
}

func (conversation *Conversation) AddBasicMessage(message textgeneration.BasicMessage) {
	conversation.Messages = append(conversation.Messages, Message{
		Role:        message.Role,
		Text:        message.Text,
		CanHaveClip: message.Role == textgeneration.RoleAssistent,
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
