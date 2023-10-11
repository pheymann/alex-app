package integrationtest_cdc

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/cmd/applogs"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
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
	generatedMessage string,
) (events.APIGatewayProxyResponse, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ctx := startartconversation.HandlerCtx{
		ConversationStore: testutil.MockConversationStore(conversations),
		UserStore:         testutil.MockUserStore(users),
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

	return ctx.AWSHandler(context.TODO(), event)
}

func MockContinueConversation(
	event events.APIGatewayProxyRequest,
	users map[string]*user.User,
	conversations map[string]*conversation.Conversation,
	generatedMessage string,
) (events.APIGatewayProxyResponse, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ctx := continueconversation.HandlerCtx{
		ConversationStore: testutil.MockConversationStore(conversations),
		UserStore:         testutil.MockUserStore(users),
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

	return ctx.AWSHandler(context.TODO(), event)
}
