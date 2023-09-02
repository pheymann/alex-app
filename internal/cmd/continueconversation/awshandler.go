package continueconversation

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/shared"
	"talktome.com/internal/talktome"
)

type HandlerCtx struct {
	Ctx talktome.Context
}

type conversationRequest struct {
	Question string `json:"question"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	convUUID := event.PathParameters["uuid"]

	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		log.Err(error).Msg("failed to extract user uuid")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to extract user uuid",
		}, nil
	}

	var convReq conversationRequest

	if err := json.Unmarshal([]byte(event.Body), &convReq); err != nil {
		log.Err(err).Msg("couldn't parse body")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Couldn't parse body",
		}, nil
	}

	message, err := Handle(userUUID, convUUID, convReq.Question, handlerCtx.Ctx)
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
