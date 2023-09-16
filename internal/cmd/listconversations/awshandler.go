package listconversations

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
	userUUID, err := shared.ExtractUserUUID(event)
	if err != nil {
		return awsutil.ReturnError(err, "failed to extract user uuid", log.With())
	}
	logCtx := log.With().Str("user_uuid", userUUID)

	convCtx := conversation.Context{
		UserUUID:          userUUID,
		LogCtx:            logCtx,
		ConversationStore: handlerCtx.ConversationStore,
		UserStore:         handlerCtx.UserStore,
	}

	conversations, err := Handle(convCtx)
	if err != nil {
		return awsutil.ReturnError(err, "failed to list all conversation", logCtx)
	}

	return awsutil.ReturnSuccessJson(conversations, logCtx)
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
