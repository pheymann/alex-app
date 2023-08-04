package listconversations

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
	"talktome.com/internal/user"
)

func Handle(userUUID string, userStorage user.StorageService, convStorage conversation.StorageService) ([]conversation.ConversationRef, error) {
	log.Info().Str("user_uuid", userUUID).Msg("fetch all conversations")

	user, err := userStorage.FindUser(userUUID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	conversations, err := convStorage.FindAllConversations(user.ConversationUUIDs)
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
