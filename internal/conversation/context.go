package conversation

import (
	"github.com/rs/zerolog"
	"talktome.com/internal/assetstore"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

type Context struct {
	UserUUID         string
	ConversationUUID string
	LogCtx           zerolog.Context

	// IO
	ConversationStore entitystore.EntityStore[Conversation]
	UserStore         entitystore.EntityStore[user.User]
	AudioClipStore    assetstore.AssetStore
	TextGen           textgeneration.TextGenerationService
	SpeechGen         speechgeneration.SpeechGenerator
}
