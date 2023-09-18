package integrationtest_cdc

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/testutil"
	"talktome.com/internal/user"
)

func Test_CDC_ConversationGetConversation(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/get_conversation.yaml",
		func(
			t *testing.T,
			responseCase CDCResponseCase[conversation.Conversation],
			event events.APIGatewayProxyRequest,
			users map[string]*user.User,
			conversations map[string]*conversation.Conversation,
		) (events.APIGatewayProxyResponse, error) {
			ctx := getconversation.HandlerCtx{
				ConversationStore: testutil.MockConversationStore(conversations),
				UserStore:         testutil.MockUserStore(users),
			}

			return ctx.AWSHandler(context.TODO(), event)
		},
	)
}
