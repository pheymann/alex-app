package talktomeartcreate

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

type ArtPiece struct {
	ArtistName string `json:"artistName"`
	ArtPiece   string `json:"artPiece"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.HTTPMethod == "POST" {
		log.Debug().Msg("POSTed art conversation creation request")

		userUUID := event.Headers["User-UUID"]

		var artPiece ArtPiece

		if err := json.Unmarshal([]byte(event.Body), &artPiece); err != nil {
			log.Err(err).Msg("couldn't parse body")
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Couldn't parse body",
			}, nil
		}

		conversation, err := Handle(userUUID, artPiece.ArtistName, artPiece.ArtPiece, handlerCtx.Ctx)
		if err != nil {
			log.Err(err).Msg("failed to create art conversation")
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to get or create conversation",
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

	log.Error().Msg("only POST requests are allowed.")
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Only POST requests are allowed.",
	}, nil
}
