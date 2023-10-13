package continueconversation

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
)

const (
	MaxMessageLength = 500
)

func Handle(ctx conversation.Context, message string) (*conversation.Message, error) {
	if message == "" {
		return nil, &shared.UserInputError{Message: "message cannot be empty"}
	} else if len(message) > MaxMessageLength {
		shared.GetLogger(ctx.LogCtx).Warn().Msgf("message too long: %d", len(message))
		return nil, &shared.UserInputError{Message: "message too long"}
	}

	conv, err := ctx.ContinueConversation(message)
	if err != nil {
		return nil, err
	}

	return &conv.Messages[len(conv.Messages)-1], nil
}
