package continueconversation

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
)

const (
	MaxMessageLength = 500
)

func Handle(ctx conversation.Context, message string) error {
	if message == "" {
		return &shared.UserInputError{Message: "message cannot be empty"}
	} else if len(message) > MaxMessageLength {
		shared.GetLogger(ctx.LogCtx).Warn().Msgf("message too long: %d", len(message))
		return &shared.UserInputError{Message: "message too long"}
	}

	// check that conversation exists
	conv, err := ctx.ConversationStore.Find(ctx.ConversationUUID, ctx.LogCtx)
	if err != nil {
		return err
	} else if conv == nil {
		return &shared.NotFoundError{Message: "conversation not found"}
	}

	if conv.State != conversation.StateReady {
		return &shared.NotFoundError{Message: "assistant is not ready yet for another question"}
	}

	// check that this user owns this conversation
	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return err
	} else if user == nil {
		return &shared.NotFoundError{Message: "user not found"}
	}

	if !user.HasConversation(conv.ID) {
		return &shared.NotFoundError{Message: "user does not have this conversation"}
	}

	conv.State = conversation.StateGenerating
	if err := ctx.ConversationStore.Save(*conv, ctx.LogCtx); err != nil {
		return err
	}

	task := processqueue.Task{
		ConversationUUID: ctx.ConversationUUID,
		UserUUID:         ctx.UserUUID,
		Language:         ctx.Language,
		Message:          message,
	}

	return ctx.ProcessQueue.Enqueue(task, ctx.LogCtx)
}
