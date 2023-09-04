package talktome

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
)

func (ctx Context) ContinueConversation(userUUID string, convUUID string, message string) (*conversation.Conversation, error) {
	log.Info().Str("user_uuid", userUUID).Str("conversation_uuid", convUUID).Msg("continue conversation")

	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

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

	clipFileName, err := ctx.conversationStorage.StoreClip(clipFile)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipUUID = clipFileName

	preSignedURL, expirationDate, err := ctx.conversationStorage.GenerateClipAccess(clipFileName)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipURL = preSignedURL
	conv.Messages[lastMessageIndex].SpeechClipExpirationDate = expirationDate
	conv.Messages[lastMessageIndex].SpeechClipIsExpired = false

	if err := ctx.conversationStorage.StoreConversation(conv); err != nil {
		return nil, err
	}

	return &conv, nil
}
