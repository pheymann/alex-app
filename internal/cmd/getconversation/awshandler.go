package getconversation

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/awsutil"
	"talktome.com/internal/conversation"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	ConversationStore entitystore.EntityStore[conversation.Conversation]
	UserStore         entitystore.EntityStore[user.User]
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	convUUID := event.PathParameters["uuid"]
	if convUUID == "" {
		return awsutil.ReturnError(nil, "conversation uuid is empty", log.With())
	}
	logCtx := log.With().Str("conversation_uuid", convUUID)

	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		return awsutil.ReturnError(error, "failed to extract user uuid", logCtx)
	}
	logCtx.Str("user_uuid", userUUID)

	convCtx := conversation.Context{
		ConversationUUID:  convUUID,
		UserUUID:          userUUID,
		LogCtx:            logCtx,
		ConversationStore: handlerCtx.ConversationStore,
		UserStore:         handlerCtx.UserStore,
	}

	conversation, err := Handle(convCtx)
	if err != nil {
		return awsutil.ReturnError(err, "failed to get conversation", logCtx)
	}

	return awsutil.ReturnSuccessJson(conversation, logCtx)
}

func UnsafeNewHandlerCtx(
	sess *session.Session,
	conversationDynamoDBTable string,
	userTable string,
) HandlerCtx {
	dynamoDBClient := dynamodb.New(sess)

	// internal init
	convStorage := entitystore.NewAWSDynamoDBCtx[conversation.Conversation](dynamoDBClient, conversationDynamoDBTable)
	userStorage := entitystore.NewAWSDynamoDBCtx[user.User](dynamoDBClient, userTable)

	return HandlerCtx{
		ConversationStore: convStorage,
		UserStore:         userStorage,
	}
}
