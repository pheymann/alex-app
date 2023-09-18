package integrationtest_cdc

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/conversation"
	"talktome.com/internal/testutil"
	"talktome.com/internal/user"
)

func Test_CDC_HomeListConversations(t *testing.T) {
	RunContracts[[]conversation.ConversationRef](
		t,
		"/home/list_conversations.yaml",
		func(
			t *testing.T,
			responseCase CDCResponseCase[[]conversation.ConversationRef],
			event events.APIGatewayProxyRequest,
			users map[string]*user.User,
			conversations map[string]*conversation.Conversation,
		) (events.APIGatewayProxyResponse, error) {
			listConvCtx := listconversations.HandlerCtx{
				ConversationStore: testutil.MockConversationStore(conversations),
				UserStore:         testutil.MockUserStore(users),
			}

			return listConvCtx.AWSHandler(context.TODO(), event)
		},
	)
}
