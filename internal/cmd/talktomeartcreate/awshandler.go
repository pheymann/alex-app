package talktomeartcreate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"talktome.com/internal/talktome"
)

type HandlerCtx struct {
	Ctx talktome.Context
}

type ArtPiece struct {
	ArtistName string `json:"artistName"`
	ArtPiece   string `json:"artPiece"`
	UserUUID   string `json:"userUUID"`
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.HTTPMethod == "POST" {
		var artPiece ArtPiece

		if err := json.Unmarshal([]byte(event.Body), &artPiece); err != nil {
			fmt.Printf("[ERROR] Couldn't parse body: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Couldn't parse body",
			}, nil
		}

		conversation, err := Handle(artPiece.UserUUID, artPiece.ArtistName, artPiece.ArtPiece, handlerCtx.Ctx)
		if err != nil {
			fmt.Printf("[ERROR] Failed to create art conversation: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to get or create conversation",
			}, nil
		}

		jsonPresentation, err := json.Marshal(*conversation)
		if err != nil {
			fmt.Printf("[ERROR] Failed tp marshal conversation: %s\n", err)
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

	fmt.Printf("[ERROR] Only POST requests are allowed.\n")
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Only POST requests are allowed.",
	}, nil
}
