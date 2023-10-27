package assistant

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
)

func Handle(ctx conversation.Context, message string) error {
	conv, err := ctx.ContinueConversation(message)
	if err != nil {
		return err
	}

	conv.State = conversation.StateReady
	if err := ctx.ConversationStore.Save(*conv, ctx.LogCtx); err != nil {
		return err
	}
	shared.GetLogger(ctx.LogCtx).Debug().Msg("assistant ready")

	return nil
}
