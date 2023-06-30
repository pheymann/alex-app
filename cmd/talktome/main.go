package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"talktome.com/internal/art"
	"talktome.com/internal/cmd/talktome"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
)

type handlerCtx struct {
	talkttome talktome.TalkToMe
}

func (handlerCtx handlerCtx) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.HTTPMethod == "POST" {
		var artPiece art.ArtPiece

		if err := json.Unmarshal([]byte(event.Body), &artPiece); err != nil {
			fmt.Printf("[ERROR] Couldn't parse body: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Couldn't parse body",
			}, nil
		}

		presentation, err := handlerCtx.talkttome.GetOrCreatePresentation(artPiece)
		if err != nil {
			fmt.Printf("[ERROR] Failed to get or create presentation: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to get or create presentation",
			}, nil
		}

		jsonPresentation, err := json.Marshal(presentation)
		if err != nil {
			fmt.Printf("[ERROR] Failed tp marshal presentation: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed tp marshal presentation",
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

func main() {
	// ENV VAR init
	openAIToken := shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	resembleToken := shared.MustReadEnvVar("TALKTOME_RESEMBLE_TOKEN")
	resembleProjectUUID := shared.MustReadEnvVar("TALKTOME_RESEMBLE_PROJECT_UUID")
	serviceDomain := shared.MustReadEnvVar("TALKTOME_SERVICE_DOMAIN")
	resembleCallBackURL := fmt.Sprintf("https://%s/callback/clip", serviceDomain)
	artPresentationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_ART_PRESENTATION_TABLE")

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	dynamoDBClient := dynamodb.New(sess)

	// internal init
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewResembleGenerator(resembleToken, resembleProjectUUID, resembleCallBackURL)
	artStorage := art.NewStorageCtx(dynamoDBClient, artPresentationDynamoDBTable, nil, "")

	talktome := talktome.NewTalkToMe(textGen, speechGen, artStorage)

	lambda.Start(handlerCtx{talkttome: talktome}.handler)
}
