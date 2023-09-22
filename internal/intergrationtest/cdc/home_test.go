package integrationtest_cdc

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/conversation"
	"talktome.com/internal/user"
)

func Test_CDC_HomeNoConversation(t *testing.T) {
	RunContracts[[]conversation.ConversationRef](
		t,
		"/home/no_conversations.yaml",
		runHomeContracts,
	)
}

func Test_CDC_HomeMultipleConversations(t *testing.T) {
	RunContracts[[]conversation.ConversationRef](
		t,
		"/home/multiple_conversations.yaml",
		runHomeContracts,
	)
}

func Test_CDC_HomeStartConversation(t *testing.T) {
	RunContracts[[]conversation.ConversationRef](
		t,
		"/home/start_conversation.yaml",
		runHomeContracts,
	)
}

func runHomeContracts(
	t *testing.T,
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
) (events.APIGatewayProxyResponse, error) {
	return MockListConversations(event, users, conversations)
}
