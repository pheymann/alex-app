package shared

import "github.com/aws/aws-lambda-go/events"

func NewAwsTestRequestContext(userUUID string) events.APIGatewayProxyRequestContext {
	return events.APIGatewayProxyRequestContext{
		Authorizer: map[string]interface{}{
			"jwt": map[string]interface{}{
				"claims": map[string]interface{}{
					"cognito:username": userUUID,
				},
			},
		},
	}
}
