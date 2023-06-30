package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"talktome.com/internal/art"
	"talktome.com/internal/cmd/resemblecallback"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
)

type handlerCtx struct {
	callback resemblecallback.ResembleCallBack
}

type resembleBody struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func (handlerCtx handlerCtx) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.HTTPMethod == "POST" {
		var body resembleBody

		if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
			fmt.Printf("[ERROR] Couldn't parse resemble body: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Couldn't parse body",
			}, nil
		}

		if err := handlerCtx.callback.StoreSpeechClip(body.ID, body.URL); err != nil {
			fmt.Printf("[ERROR] Failed to download and store clip: %s\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Failed to download and store clip",
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
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
	resembleToken := shared.MustReadEnvVar("TALKTOME_RESEMBLE_TOKEN")
	resembleProjectUUID := shared.MustReadEnvVar("TALKTOME_RESEMBLE_PROJECT_UUID")
	artSpeechClipBucket := shared.MustReadEnvVar("TALKTOME_ART_SPEECH_CLIP_BUCKET")

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your desired region
	})
	if err != nil {
		panic(err)
	}

	s3Client := s3.New(sess)

	// internal init
	speechGen := speechgeneration.NewResembleGenerator(resembleToken, resembleProjectUUID, "")
	artStorage := art.NewStorageCtx(nil, "", s3Client, artSpeechClipBucket)
	callback := resemblecallback.NewResembleCallBack(speechGen, artStorage)

	lambda.Start(handlerCtx{callback: callback}.handler)
}
