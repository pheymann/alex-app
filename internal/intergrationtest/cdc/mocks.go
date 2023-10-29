package integrationtest_cdc

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/cmd/applogs"
	"talktome.com/internal/cmd/assistant"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/pollassistantresponse"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/testutil"
	"talktome.com/internal/user"
)

func MockGetConversation(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
) (events.APIGatewayProxyResponse, error) {
	ctx := getconversation.HandlerCtx{
		ConversationStore: testutil.MockConversationStore(conversations),
		UserStore:         testutil.MockUserStore(users),
		AudioClipStore: &testutil.MockAssetStore{
			PresignedUrl: "https://some.url/clip0",
		},
	}

	return ctx.AWSHandler(context.TODO(), event)
}

func MockAppLogs(
	event events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return applogs.AWSHandler(context.TODO(), event)
}

func MockListConversations(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
) (events.APIGatewayProxyResponse, error) {
	ctx := listconversations.HandlerCtx{
		ConversationStore: testutil.MockConversationStore(conversations),
		UserStore:         testutil.MockUserStore(users),
	}

	return ctx.AWSHandler(context.TODO(), event)
}

func MockStartArtConversations(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
	startedConvID string,
	generatedMessage string,
) (events.APIGatewayProxyResponse, error) {
	mockConvStore := testutil.MockConversationStore(conversations)
	mockUserStore := testutil.MockUserStore(users)
	mockProcessQueue := testutil.MockProcessQueue{
		Queue: make(chan processqueue.Task, 1),
	}

	ctx := startartconversation.HandlerCtx{
		ConversationStore: mockConvStore,
		UserStore:         mockUserStore,
		ProcessQueue:      &mockProcessQueue,
		IDGenerator: &testutil.MockIDGenerator{
			GeneratedID: startedConvID,
			UseMetadata: false,
		},
	}

	response, err := ctx.AWSHandler(context.TODO(), event)
	if err != nil || response.StatusCode > 300 {
		return response, err
	}

	// trigger async assistant to make sure an answer is
	// ready with the next poll
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	assistantCtx := assistant.HandlerCtx{
		ConversationStore: mockConvStore,
		UserStore:         mockUserStore,
		AudioClipStore: &testutil.MockAssetStore{
			ClipKey:      "clip0",
			PresignedUrl: "https://some.url/clip0",
		},
		TextGen: &testutil.MockTextGeneration{
			GeneratedMessage: generatedMessage,
			Timeout:          0 * time.Second,
		},
		SpeechGen: &testutil.MockSpeechGeneration{
			TestFile: rootPath + "/../../../assets/prompt.mp3",
			Timeout:  0 * time.Second,
		},
	}

	task := <-mockProcessQueue.Queue
	taskJson, err := json.Marshal(task)
	if err != nil {
		return response, err
	}

	err = assistantCtx.AWSHandler(context.TODO(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(taskJson),
			},
		},
	})
	if err != nil {
		return response, err
	}
	// done with async processing

	return response, nil
}

func MockContinueConversation(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
	generatedMessage string,
) (events.APIGatewayProxyResponse, error) {
	mockConvStore := testutil.MockConversationStore(conversations)
	mockUserStore := testutil.MockUserStore(users)
	mockProcessQueue := testutil.MockProcessQueue{
		Queue: make(chan processqueue.Task, 1),
	}

	ctx := continueconversation.HandlerCtx{
		ConversationStore: mockConvStore,
		UserStore:         mockUserStore,
		ProcessQueue:      &mockProcessQueue,
	}

	response, err := ctx.AWSHandler(context.TODO(), event)
	if err != nil || response.StatusCode > 300 {
		return response, err
	}

	// trigger async assistant to make sure an answer is
	// ready with the next poll
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	assistantCtx := assistant.HandlerCtx{
		ConversationStore: mockConvStore,
		UserStore:         mockUserStore,
		AudioClipStore: &testutil.MockAssetStore{
			ClipKey:      "clip1",
			PresignedUrl: "https://some.url/clip1",
		},
		TextGen: &testutil.MockTextGeneration{
			GeneratedMessage: generatedMessage,
			Timeout:          0 * time.Second,
		},
		SpeechGen: &testutil.MockSpeechGeneration{
			TestFile: rootPath + "/../../../assets/prompt.mp3",
			Timeout:  0 * time.Second,
		},
	}

	task := <-mockProcessQueue.Queue
	taskJson, err := json.Marshal(task)
	if err != nil {
		return response, err
	}

	err = assistantCtx.AWSHandler(context.TODO(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(taskJson),
			},
		},
	})
	if err != nil {
		return response, err
	}
	// done with async processing

	return response, nil
}

func MockPollAssistantResponse(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
) (events.APIGatewayProxyResponse, error) {
	ctx := pollassistantresponse.HandlerCtx{
		ConversationStore: testutil.MockConversationStore(conversations),
		UserStore:         testutil.MockUserStore(users),
	}

	return ctx.AWSHandler(context.TODO(), event)
}
