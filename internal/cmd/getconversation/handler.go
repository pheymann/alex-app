package getconversation

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
)

func Handle(ctx conversation.Context) (*conversation.Conversation, error) {
	shared.GetLogger(ctx.LogCtx).Debug().Msg("getting conversation")

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

	if user.HasConversation(ctx.ConversationUUID) {
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
			if message.Role == textgeneration.RoleAssistent {
				preSignedURL, err := ctx.AudioClipStore.GenerateTemporaryAccessURL(message.SpeechClipUUID, ctx.LogCtx)
				if err != nil {
					return nil, err
				}

				conv.Messages[index].SpeechClipURL = preSignedURL
			}
		}

		return conv, nil
	}

	return nil, &shared.NotFoundError{Message: "user does not have this conversation"}
}
