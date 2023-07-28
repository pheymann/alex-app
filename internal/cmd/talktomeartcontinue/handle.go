package talktomeartcontinue

import (
	"talktome.com/internal/conversation"
	"talktome.com/internal/talktome"
)

func Handle(userUUID string, convUUID string, message string, ctx talktome.Context) (*conversation.Conversation, error) {
	return ctx.ContinueConversation(userUUID, convUUID, message)
}
