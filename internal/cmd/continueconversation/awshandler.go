package continueconversation

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/talktome"
)

type HandlerCtx struct {
	Ctx talktome.Context
}

type conversationRequest struct {
	ConversationUUID string `json:"conversationUuid"`
	Prompt           string `json:"prompt"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.HTTPMethod == "POST" {
		log.Debug().Msg("POSTed conversation continuation request")
		// TODO: add user id to log context
		userUUID := event.Headers["User-UUID"]

		var convReq conversationRequest

		if err := json.Unmarshal([]byte(event.Body), &convReq); err != nil {
			log.Err(err).Msg("couldn't parse body")
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Couldn't parse body",
			}, nil
		}

		message, err := Handle(userUUID, convReq.ConversationUUID, convReq.Prompt, handlerCtx.Ctx)
		if err != nil {
			log.Err(err).Msg("failed to continue conversation")
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to continue conversation",
			}, nil
		}

		jsonPresentation, err := json.Marshal(*message)
		if err != nil {
			log.Err(err).Msg("failed to marshal response")
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to marshal response",
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(jsonPresentation),
		}, nil
	}

	log.Error().Msg("only POST requests are allowed.")
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Only POST requests are allowed.",
	}, nil
}