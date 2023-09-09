package conversation

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	ID       string            `json:"id" dynamodbav:"id"`
	Metadata map[string]string `json:"metadata" dynamodbav:"metadata"`
	Messages []Message         `json:"messages" dynamodbav:"messages"`
}

type Message struct {
	Role                     Role       `json:"role" dynamodbav:"role"`
	Text                     string     `json:"text" dynamodbav:"text"`
	CanHaveClip              bool       `json:"canHaveClip" dynamodbav:"can_have_clip"`
	SpeechClipUUID           string     `json:"speechClipUuid" dynamodbav:"speech_clip_uuid"`
	SpeechClipURL            string     `json:"speechClipUrl" dynamodbav:"speech_clip_url"`
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
		Role:        openai.ChatMessageRoleUser,
		Text:        text,
		CanHaveClip: false,
	})
}
