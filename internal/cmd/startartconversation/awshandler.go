package startartconversation

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/assetstore"
	"talktome.com/internal/cmd/awsutil"
	"talktome.com/internal/conversation"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	ConversationStore entitystore.EntityStore[conversation.Conversation]
	UserStore         entitystore.EntityStore[user.User]
	AudioClipStore    assetstore.AssetStore
	TextGen           textgeneration.TextGenerationService
	SpeechGen         speechgeneration.SpeechGenerator
}

type ArtContext struct {
	Context string `json:"artContext"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		return awsutil.ReturnError(error, "failed to extract user uuid", log.With())
	}
	logCtx := log.With().Str("user_uuid", userUUID)

	var artContext ArtContext

	if err := json.Unmarshal([]byte(event.Body), &artContext); err != nil {
		return awsutil.ReturnError(err, "couldn't parse body", logCtx)
	}

	convCtx := conversation.Context{
		UserUUID:          userUUID,
		LogCtx:            logCtx,
		ConversationStore: handlerCtx.ConversationStore,
		UserStore:         handlerCtx.UserStore,
		AudioClipStore:    handlerCtx.AudioClipStore,
		TextGen:           handlerCtx.TextGen,
		SpeechGen:         handlerCtx.SpeechGen,
	}

	conversation, err := Handle(convCtx, artContext.Context)
	if err != nil {
		return awsutil.ReturnError(err, "failed to start art conversation", logCtx)
	}

	return awsutil.ReturnSuccessJson(conversation, logCtx)
}

func UnsafeNewHandlerCtx(
	sess *session.Session,
	conversationDynamoDBTable string,
	userTable string,
	openAIToken string,
	conversationClipBucket string,
) HandlerCtx {
	dynamoDBClient := dynamodb.New(sess)
	s3 := s3.New(sess)
	pollyClient := polly.New(sess)

	// internal init
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewAWSPollySpeechGenerator(pollyClient)
	convStorage := entitystore.NewAWSDynamoDBCtx[conversation.Conversation](dynamoDBClient, conversationDynamoDBTable)
	userStorage := entitystore.NewAWSDynamoDBCtx[user.User](dynamoDBClient, userTable)
	audioClipStore := assetstore.NewAWSS3Context(s3, conversationClipBucket)

	return HandlerCtx{
		ConversationStore: convStorage,
		UserStore:         userStorage,
		AudioClipStore:    audioClipStore,
		TextGen:           textGen,
		SpeechGen:         speechGen,
	}
}
