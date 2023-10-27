package startartconversation

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/awsutil"
	"talktome.com/internal/conversation"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/idgenerator"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	ConversationStore entitystore.EntityStore[conversation.Conversation]
	UserStore         entitystore.EntityStore[user.User]
	ProcessQueue      processqueue.ProcessQueue
	IDGenerator       idgenerator.IDGenerator
}

type ArtContext struct {
	Context string `json:"artContext"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		return awsutil.ReturnError(error, "failed to extract user uuid", log.With())
	}
	logCtx := log.With().Str("lambda", "startartconversation").Str("user_uuid", userUUID)

	var language = shared.LanguageGerman
	if languageStr, ok := event.Headers["accept-language"]; ok {
		shared.GetLogger(logCtx).Debug().Msgf("found language header: %s", languageStr)
		language = shared.DecodeLanguage(languageStr)
	}
	logCtx = logCtx.Str("language", string(language))

	var artContext ArtContext

	if err := json.Unmarshal([]byte(event.Body), &artContext); err != nil {
		return awsutil.ReturnError(err, "couldn't parse body", logCtx)
	}

	convCtx := conversation.Context{
		UserUUID:          userUUID,
		Language:          language,
		LogCtx:            logCtx,
		ConversationStore: handlerCtx.ConversationStore,
		UserStore:         handlerCtx.UserStore,
		ProcessQueue:      handlerCtx.ProcessQueue,
		IDGenerator:       handlerCtx.IDGenerator,
	}

	conv, err := Handle(convCtx, artContext.Context)
	if err != nil {
		return awsutil.ReturnError(err, "failed to start art conversation", logCtx)
	}

	return awsutil.ReturnSuccessJson(conv, logCtx)
}

func UnsafeNewHandlerCtx(
	sess *session.Session,
	conversationDynamoDBTable string,
	userTable string,
	queueURL string,
) HandlerCtx {
	dynamoDBClient := dynamodb.New(sess)
	sqs := sqs.New(sess)

	// internal init
	convStorage := entitystore.NewAWSDynamoDBCtx[conversation.Conversation](dynamoDBClient, conversationDynamoDBTable)
	userStorage := entitystore.NewAWSDynamoDBCtx[user.User](dynamoDBClient, userTable)
	processQueue := processqueue.NewAWSSQSContext(sqs, queueURL)

	return HandlerCtx{
		ConversationStore: convStorage,
		UserStore:         userStorage,
		ProcessQueue:      processQueue,
		IDGenerator:       idgenerator.NewRandomIDGenerator(),
	}
}
