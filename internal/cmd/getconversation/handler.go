package getconversation

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
	"talktome.com/internal/user"
)

func Handle(userUUID string, convUUID string, userStorage user.StorageService, convStorage conversation.StorageService) (*conversation.Conversation, error) {
	log.Info().Str("user_uuid", userUUID).Str("conv_uuid", convUUID).Msg("get conversation")

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("failed to load UTC location: %w", err)
	}

	user, err := userStorage.FindUser(userUUID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	conv, err := convStorage.FindConversation(convUUID)
	if err != nil {
		return nil, err
	} else if conv == nil {
		return nil, fmt.Errorf("conversation not found")
	}

	for _, id := range user.ConversationUUIDs {
		if id == convUUID {
			conv.Messages = conv.Messages[3:]

			for index, message := range conv.Messages {
				if message.SpeechClipExpirationDate != nil {
					isExpired := time.Now().In(location).After(*message.SpeechClipExpirationDate)
					conv.Messages[index].SpeechClipIsExpired = isExpired
				}
			}

			return conv, nil
		}
	}

	return nil, fmt.Errorf("user does not have this conversation")
}
