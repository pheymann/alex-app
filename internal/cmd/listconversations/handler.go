package listconversations

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
)

func Handle(ctx conversation.Context) ([]conversation.ConversationRef, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msg("fetch all conversations")

	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, &shared.NotFoundError{Message: "user not found"}
	}

	if len(user.ConversationUUIDs) == 0 {
		return []conversation.ConversationRef{}, nil
	}

	conversations, err := ctx.ConversationStore.FindAll(user.ConversationUUIDs, ctx.LogCtx)
	if err != nil {
		return nil, err
	}

	conversationRef := make([]conversation.ConversationRef, len(conversations))
	for index, conv := range conversations {
		conversationRef[index] = conversation.ConversationRef{
			ID:       conv.ID,
			Metadata: conv.Metadata,
		}
	}

	return conversationRef, nil
}
