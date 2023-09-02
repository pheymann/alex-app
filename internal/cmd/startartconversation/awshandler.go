package startartconversation

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

type ArtContext struct {
	Context string `json:"artContext"`
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

	var artContext ArtContext

	if err := json.Unmarshal([]byte(event.Body), &artContext); err != nil {
		log.Err(err).Msg("couldn't parse body")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Couldn't parse body",
		}, nil
	}

	conversation, err := Handle(userUUID, artContext.Context, handlerCtx.Ctx)
	if err != nil {
		log.Err(err).Msg("failed to start art conversation")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to get or start conversation",
		}, nil
	}

	jsonPresentation, err := json.Marshal(*conversation)
	if err != nil {
		log.Err(err).Msg("failed to marshal conversation")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed tp marshal conversation",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonPresentation),
	}, nil
}
