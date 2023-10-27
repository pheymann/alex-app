package processqueue

import (
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

type Task struct {
	ConversationUUID string          `json:"conversation_uuid"`
	UserUUID         string          `json:"user_uuid"`
	Language         shared.Language `json:"language"`
	Message          string          `json:"message"`
}

type ProcessQueue interface {
	Enqueue(task Task, logCtx zerolog.Context) error
}
