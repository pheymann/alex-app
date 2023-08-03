package listconversations

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
	"talktome.com/internal/user"
)

type ConversationRef struct {
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata"`
}

func Handle(userUUID string, userStorage user.StorageService, convStorage conversation.StorageService) ([]ConversationRef, error) {
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

	conversationRef := make([]ConversationRef, len(conversations))
	for index, conversation := range conversations {
		conversationRef[index] = ConversationRef{
			ID:       conversation.ID,
			Metadata: conversation.Metadata,
		}
	}

	return conversationRef, nil
}
