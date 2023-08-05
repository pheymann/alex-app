package talktome

import (
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
)

func (ctx Context) ContinueConversation(userUUID string, convUUID string, message string) (*conversation.Conversation, error) {
	log.Info().Str("user_uuid", userUUID).Str("conversation_uuid", convUUID).Msg("continue conversation")

	user, err := ctx.userStorage.FindUser(userUUID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user with UUID %s does not exist", userUUID)
	}

	optConv, err := ctx.conversationStorage.FindConversation(convUUID)
	if err != nil {
		return nil, err
	} else if optConv == nil {
		return nil, fmt.Errorf("conversation with UUID %s does not exist", convUUID)
	}

	conv := *optConv
	conv.AddMessage(message)

	if err := ctx.textGen.GenerateNextMessage(&conv); err != nil {
		return nil, err
	}

	lastMessageIndex := len(conv.Messages) - 1

	clipFile, err := ctx.speechGen.GenerateSpeechClip(conv.ID, conv.Messages[lastMessageIndex].Text)
	if err != nil {
		return nil, err
	}
	defer clipFile.Close()

	if err := ctx.conversationStorage.StoreClip(clipFile); err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipUUID = filepath.Base(clipFile.Name())

	if err := ctx.conversationStorage.StoreConversation(conv); err != nil {
		return nil, err
	}

	return &conv, nil
}
