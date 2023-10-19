package continueconversation

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

type conversationRequest struct {
	Question string `json:"question"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	convUUID := event.PathParameters["uuid"]
	if convUUID == "" {
		return awsutil.ReturnError(nil, "conversation uuid is empty", log.With())
	}
	var logCtx = log.With().Str("conversation_uuid", convUUID)

	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		return awsutil.ReturnError(error, "failed to extract user uuid", logCtx)
	}
	logCtx = logCtx.Str("user_uuid", userUUID)

	var language = shared.LanguageGerman
	if languageStr, ok := event.Headers["Accept-Language"]; ok {
		shared.GetLogger(logCtx).Debug().Msgf("found language header: %s", languageStr)
		language = shared.DecodeLanguage(languageStr)
	}
	logCtx = logCtx.Str("language", string(language))

	shared.GetLogger(logCtx).Info().Msg("continuing conversation")

	var convReq conversationRequest

	if err := json.Unmarshal([]byte(event.Body), &convReq); err != nil {
		return awsutil.ReturnError(err, "couldn't parse body", logCtx)
	}

	convCtx := conversation.Context{
		ConversationUUID:  convUUID,
		UserUUID:          userUUID,
		Language:          language,
		LogCtx:            logCtx,
		ConversationStore: handlerCtx.ConversationStore,
		UserStore:         handlerCtx.UserStore,
		AudioClipStore:    handlerCtx.AudioClipStore,
		TextGen:           handlerCtx.TextGen,
		SpeechGen:         handlerCtx.SpeechGen,
	}

	message, err := Handle(convCtx, convReq.Question)
	if err != nil {
		return awsutil.ReturnError(err, "failed to continue conversation", logCtx)
	}

	return awsutil.ReturnSuccessJson(message, logCtx)
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
