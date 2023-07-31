package talktome

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
)

func (ctx Context) StartConversation(userUUID string, conv conversation.Conversation, message string) (*conversation.Conversation, error) {
	log.Info().Str("user_uuid", userUUID).Str("conversation_uuid", conv.ID).Msg("start conversation")

	user, err := ctx.userStorage.FindUser(userUUID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user with UUID %s does not exist", userUUID)
	}

	if err := ctx.conversationStorage.StoreConversation(conv); err != nil {
		return nil, err
	}

	user.ConversationUUIDs = append(user.ConversationUUIDs, conv.ID)
	if err := ctx.userStorage.StoreUser(*user); err != nil {
		return nil, err
	}

	continuedConv, err := ctx.ContinueConversation(userUUID, conv.ID, message)
	if err != nil {
		return nil, err
	}

	continuedConv.Messages = continuedConv.Messages[3:]
	return continuedConv, nil
}
