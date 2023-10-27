package conversation

import (
	"github.com/rs/zerolog"
	"talktome.com/internal/assetstore"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/idgenerator"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

type Context struct {
	UserUUID         string
	ConversationUUID string
	Language         shared.Language
	LogCtx           zerolog.Context

	// IO
	ConversationStore entitystore.EntityStore[Conversation]
	UserStore         entitystore.EntityStore[user.User]
	AudioClipStore    assetstore.AssetStore
	TextGen           textgeneration.TextGenerationService
	SpeechGen         speechgeneration.SpeechGenerator
	ProcessQueue      processqueue.ProcessQueue
	IDGenerator       idgenerator.IDGenerator
}
