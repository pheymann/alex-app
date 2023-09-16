package conversation

import (
	"fmt"

	"talktome.com/internal/shared"
)

func (ctx Context) ContinueConversation(message string) (*Conversation, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msgf("continue conversation with message: %s", message)

	if message == "" {
		return nil, &shared.UserInputError{Message: "message cannot be empty"}
	}

	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, &shared.NotFoundError{Message: fmt.Sprintf("user with UUID %s does not exist", ctx.UserUUID)}
	}

	optConv, err := ctx.ConversationStore.Find(ctx.ConversationUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if optConv == nil {
		return nil, &shared.NotFoundError{Message: fmt.Sprintf("conversation with UUID %s does not exist", ctx.ConversationUUID)}
	}

	conv := *optConv
	conv.AddMessage(message)

	nextMessage, err := ctx.TextGen.GenerateNextMessage(conv.ToBasicMessageHistory(), ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conv.AddBasicMessage(*nextMessage)

	lastMessageIndex := len(conv.Messages) - 1

	clipFile, err := ctx.SpeechGen.GenerateSpeechClip(conv.ID, conv.Messages[lastMessageIndex].Text, ctx.LogCtx)
	if err != nil {
		return nil, err
	}
	defer clipFile.Close()

	clipFileName, err := ctx.AudioClipStore.Save(clipFile, ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipUUID = clipFileName

	preSignedURL, expirationDate, err := ctx.AudioClipStore.GenerateTemporaryAccessURL(clipFileName, ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipURL = preSignedURL
	conv.Messages[lastMessageIndex].SpeechClipExpirationDate = expirationDate
	conv.Messages[lastMessageIndex].SpeechClipIsExpired = false

	if err := ctx.ConversationStore.Save(conv, ctx.LogCtx); err != nil {
		return nil, err
	}

	return &conv, nil
}
