package conversation

import (
	"fmt"

	"talktome.com/internal/shared"
)

func (ctx Context) ContinueConversation(message string) (*Conversation, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msgf("continue conversation with message: %s", message)

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

	if err := ctx.ConversationStore.Save(conv, ctx.LogCtx); err != nil {
		return nil, err
	}

	lastMessageIndex := len(conv.Messages) - 1

	clipFile, err := ctx.SpeechGen.GenerateSpeechClip(conv.ID, conv.Messages[lastMessageIndex].Text, ctx.Language, ctx.LogCtx)
	if err != nil {
		return nil, err
	}
	defer clipFile.Close()

	clipFileName, err := ctx.AudioClipStore.Save(clipFile, ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipUUID = clipFileName

	preSignedURL, err := ctx.AudioClipStore.GenerateTemporaryAccessURL(clipFileName, ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conv.Messages[lastMessageIndex].SpeechClipURL = preSignedURL

	if err := ctx.ConversationStore.Save(conv, ctx.LogCtx); err != nil {
		return nil, err
	}

	return &conv, nil
}
