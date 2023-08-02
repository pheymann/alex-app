package conversation

import (
	"encoding/base64"

	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Information string    `json:"information" dynamodbav:"information"`
	Messages    []Message `json:"messages" dynamodbav:"messages"`
}

type Message struct {
	Role           Role   `json:"role" dynamodbav:"role"`
	Text           string `json:"text" dynamodbav:"text"`
	CanHaveClip    bool   `json:"canHaveClip" dynamodbav:"can_have_clip"`
	SpeechClipUUID string `json:"speechClipUuid" dynamodbav:"speech_clip_uuid"`
}

type Role = string

const (
	RoleUser      Role = "user"
	RoleAssistent Role = "assistant"
	RoleSystem    Role = "system"
)

func NewConversation(information string) Conversation {
	return Conversation{
		ID:          GenerateStableID(information),
		Information: information,
		Messages:    []Message{},
	}
}

func GenerateStableID(information string) string {
	return base64.StdEncoding.EncodeToString([]byte(information))
}

func (conversation *Conversation) AddMessage(text string) {
	conversation.Messages = append(conversation.Messages, Message{
		Role:        openai.ChatMessageRoleUser,
		Text:        text,
		CanHaveClip: false,
	})
}
