package getconversation

import (
	"time"

	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
)

func Handle(ctx conversation.Context) (*conversation.Conversation, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msg("getting conversation")

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, &shared.InternalError{Cause: err, Message: "failed to load UTC location"}
	}

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

	for _, id := range user.ConversationUUIDs {
		if id == ctx.ConversationUUID {
			// don't show system messages (assumption: only at the beginning)
			for index, message := range conv.Messages {
				if message.Role != textgeneration.RoleSystem {
					conv.Messages = conv.Messages[index:]
					break
				}
			}

			// also ignore first question because that is generated
			conv.Messages = conv.Messages[1:]

			for index, message := range conv.Messages {
				if message.SpeechClipExpirationDate != nil {
					isExpired := time.Now().In(location).After(*message.SpeechClipExpirationDate)
					conv.Messages[index].SpeechClipIsExpired = isExpired
				}
			}

			return conv, nil
		}
	}

	return nil, &shared.NotFoundError{Message: "user does not have this conversation"}
}
