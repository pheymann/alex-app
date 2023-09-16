package continueconversation

import (
	"talktome.com/internal/conversation"
)

func Handle(ctx conversation.Context, message string) (*conversation.Message, error) {
	conv, err := ctx.ContinueConversation(message)
	if err != nil {
		return nil, err
	}

	return &conv.Messages[len(conv.Messages)-1], nil
}
