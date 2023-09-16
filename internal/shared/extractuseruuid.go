package shared

import (
	"github.com/aws/aws-lambda-go/events"
)

func ExtractUserUUID(event events.APIGatewayProxyRequest) (string, error) {
	if jwt, ok := event.RequestContext.Authorizer["jwt"].(map[string]interface{}); ok {
		if claims, ok := jwt["claims"].(map[string]interface{}); ok {
			if userUUID, ok := claims["cognito:username"].(string); ok {
				return userUUID, nil
			}
			return "", &AuthorizationError{nil, "failed to extract user uuid from claims"}
		}
		return "", &AuthorizationError{nil, "failed to extract claims from jwt"}
	}
	return "", &AuthorizationError{nil, "failed to extract jwt from request context"}
}
