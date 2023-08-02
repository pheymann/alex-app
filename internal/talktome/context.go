package talktome

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/user"
)

type Context struct {
	textGen             conversation.TextGenerationService
	speechGen           speechgeneration.SpeechGenerator
	conversationStorage conversation.StorageService
	userStorage         user.StorageService
}

func NewContext(
	textGen conversation.TextGenerationService,
	speechGen speechgeneration.SpeechGenerator,
	storage conversation.StorageService,
	userStorage user.StorageService,
) Context {
	return Context{
		textGen:             textGen,
		speechGen:           speechGen,
		conversationStorage: storage,
		userStorage:         userStorage,
	}
}
