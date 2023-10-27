package integrationtest_cdc

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/conversation"
	"talktome.com/internal/user"
)

func Test_CDC_ConversationFound(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/found.yaml",
		runConversationContracts,
	)
}

func Test_CDC_ConversationNotFound(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/not_found.yaml",
		runConversationContracts,
	)
}

func Test_CDC_ConversationUnauthorizedAccess(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/unauthorized_access.yaml",
		runConversationContracts,
	)
}

func Test_CDC_ConversationStartConversation(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/start_conversation.yaml",
		runConversationContracts,
	)
}

func Test_CDC_ConversationAskQuestion(t *testing.T) {
	RunContracts[conversation.Conversation](
		t,
		"/conversation/ask_question.yaml",
		runConversationContracts,
	)
}

func runConversationContracts(
	t *testing.T,
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
) (events.APIGatewayProxyResponse, error) {
	if strings.Contains(event.Path, "/api/conversation/list") {
		return MockListConversations(event, users, conversations)
	} else if strings.Contains(event.Path, "/api/conversation/create/art") {
		response, err := MockStartArtConversations(event, users, conversations, "abc", "some answer")
		if err != nil {
			return response, err
		}

		// make sure random IDs are set to what is in the contract
		var conversation = conversation.Conversation{}
		if err := json.Unmarshal([]byte(response.Body), &conversation); err != nil {
			t.Fatalf("failed to unmarshal conversation: %s", err)
		}
		conversation.ID = "abc"

		body, err := json.Marshal(conversation)
		if err != nil {
			t.Fatalf("failed to marshal conversation: %s", err)
		}

		response.Body = string(body)

		return response, err
	} else if strings.Contains(event.Path, "/api/conversation/") {
		if strings.Contains(event.Path, "continue") {
			return MockContinueConversation(event, users, conversations, "another answer")
		} else if strings.Contains(event.Path, "poll") {
			return MockPollAssistantResponse(event, users, conversations)
		}
		return MockGetConversation(event, users, conversations)
	} else if strings.Contains(event.Path, "/api/app/logs") {
		return MockAppLogs(event)
	} else {
		t.Fatalf("unknown path: %s", event.Path)
		return events.APIGatewayProxyResponse{}, nil
	}
}
