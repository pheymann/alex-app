package conversation

import (
	"fmt"

	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
)

func (ctx Context) StartConversation(conv Conversation, message string) (*Conversation, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msg("start conversation")

	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, &shared.NotFoundError{Message: fmt.Sprintf("user with UUID %s does not exist", ctx.UserUUID)}
	}

	if err := ctx.ConversationStore.Save(conv, ctx.LogCtx); err != nil {
		return nil, err
	}

	user.ConversationUUIDs = append(user.ConversationUUIDs, conv.ID)
	if err := ctx.UserStore.Save(*user, ctx.LogCtx); err != nil {
		return nil, err
	}

	continuedConv, err := ctx.ContinueConversation(message)
	if err != nil {
		return nil, err
	}

	// don't show system messages (assumption: only at the beginning)
	for index, message := range continuedConv.Messages {
		if message.Role != textgeneration.RoleSystem {
			continuedConv.Messages = continuedConv.Messages[index:]
			break
		}
	}

	// also ignore first question because that is generated
	continuedConv.Messages = continuedConv.Messages[1:]
	return continuedConv, nil
}
