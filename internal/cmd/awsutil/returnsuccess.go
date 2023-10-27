package awsutil

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

func ReturnSuccessJson(body any, logCtx zerolog.Context) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ReturnError(err, "failed to marshal body", logCtx)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}

func ReturnStatus(statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
	}, nil
}
