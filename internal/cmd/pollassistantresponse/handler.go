package pollassistantresponse

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
)

func Handle(ctx conversation.Context) (*conversation.Message, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msg("polling conversation")

	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, &shared.NotFoundError{Message: "user not found"}
	}

	conv, err := ctx.ConversationStore.Find(ctx.ConversationUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if conv == nil {
		return nil, &shared.NotFoundError{Message: "conversation not found"}
	}

	if !user.HasConversation(conv.ID) {
		return nil, &shared.NotFoundError{Message: "user does not have this conversation"}
	}

	if conv.State == conversation.StateReady {
		assistantResponse := conv.Messages[len(conv.Messages)-1]

		return &assistantResponse, nil
	}

	return nil, nil
}
