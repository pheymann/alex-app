package textgeneration

import "github.com/rs/zerolog"

type TextGenerationService interface {
	GenerateNextMessage(messageHistory []BasicMessage, logCtx zerolog.Context) (*BasicMessage, error)
}

type BasicMessage struct {
	Role MessageAuthorRole
	Text string
}

type MessageAuthorRole = string

const (
	RoleAssistent MessageAuthorRole = "assistant"
	RoleUser      MessageAuthorRole = "user"
	RoleSystem    MessageAuthorRole = "system"
)
