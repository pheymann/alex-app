package continueconversation

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/talktome"
)

func Handle(userUUID string, convUUID string, message string, ctx talktome.Context) (*conversation.Message, error) {
	conv, err := ctx.ContinueConversation(userUUID, convUUID, message)
	if err != nil {
		return nil, err
	}

	return &conv.Messages[len(conv.Messages)-1], nil
}
