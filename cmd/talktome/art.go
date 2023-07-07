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
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
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

		conversation, _, err := handlerCtx.talkttome.TalkToMeArt(artPiece, nil)
		if err != nil {
			fmt.Printf("[ERROR] Failed to get or create presentation: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Failed to get or create presentation",
			}, nil
		}

		jsonPresentation, err := json.Marshal(*conversation)
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
	artPresentationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_ART_PRESENTATION_TABLE")

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	dynamoDBClient := dynamodb.New(sess)
	s3 := s3.New(sess)
	pollyClient := polly.New(sess)

	// internal init
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewAWSPollySpeechGenerator(pollyClient)
	artStorage := art.NewAWSStorageCtx(dynamoDBClient, artPresentationDynamoDBTable, s3, "talktome-artaudioclips")

	talktome := talktome.NewTalkToMe(textGen, speechGen, artStorage)

	lambda.Start(handlerCtx{talkttome: talktome}.handler)
}
