package listconversations

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	UserStorage user.StorageService
	ConvStorage conversation.StorageService
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		log.Err(error).Msg("failed to extract user uuid")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to extract user uuid",
		}, nil
	}

	conversations, err := Handle(userUUID, handlerCtx.UserStorage, handlerCtx.ConvStorage)
	if err != nil {
		log.Err(err).Msg("failed to list conversations")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to list conversations",
		}, nil
	}
	log.Debug().Msgf("found %d conversations", len(conversations))

	jsonConversations, err := json.Marshal(conversations)
	if err != nil {
		log.Err(err).Msg("failed to marshal conversations")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed tp marshal conversations",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonConversations),
	}, nil
}
