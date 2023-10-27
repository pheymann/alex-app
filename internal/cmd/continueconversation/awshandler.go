package continueconversation

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
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	ConversationStore entitystore.EntityStore[conversation.Conversation]
	UserStore         entitystore.EntityStore[user.User]
	ProcessQueue      processqueue.ProcessQueue
}

type conversationRequest struct {
	Question string `json:"question"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	convUUID := event.PathParameters["uuid"]
	if convUUID == "" {
		return awsutil.ReturnError(nil, "conversation uuid is empty", log.With())
	}
	var logCtx = log.With().Str("lambda", "continueconversation").Str("conversation_uuid", convUUID)

	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		return awsutil.ReturnError(error, "failed to extract user uuid", logCtx)
	}
	logCtx = logCtx.Str("user_uuid", userUUID)

	var language = shared.LanguageGerman
	if languageStr, ok := event.Headers["accept-language"]; ok {
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
		ProcessQueue:      handlerCtx.ProcessQueue,
	}

	if err := Handle(convCtx, convReq.Question); err != nil {
		return awsutil.ReturnError(err, "failed to continue conversation", logCtx)
	}

	return awsutil.ReturnStatus(202)
}

func UnsafeNewHandlerCtx(
	sess *session.Session,
	conversationDynamoDBTable string,
	userTable string,
	QueueURL string,
) HandlerCtx {
	dynamoDBClient := dynamodb.New(sess)
	sqs := sqs.New(sess)

	// internal init
	convStorage := entitystore.NewAWSDynamoDBCtx[conversation.Conversation](dynamoDBClient, conversationDynamoDBTable)
	userStorage := entitystore.NewAWSDynamoDBCtx[user.User](dynamoDBClient, userTable)
	processQueue := processqueue.NewAWSSQSContext(sqs, QueueURL)

	return HandlerCtx{
		ConversationStore: convStorage,
		UserStore:         userStorage,
		ProcessQueue:      processQueue,
	}
}
