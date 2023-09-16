package awsutil

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"talktome.com/internal/assetstore"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/shared"
)

func ReturnError(err error, message string, logCtx zerolog.Context) (events.APIGatewayProxyResponse, error) {
	switch err.(type) {
	case *shared.AuthorizationError:
		return returnErrorWithCode(401, err, message, logCtx)

	case *shared.InternalError:
		return returnErrorWithCode(500, err, message, logCtx)

	case *shared.NotFoundError:
		return returnErrorWithCode(404, err, message, logCtx)

	case *assetstore.AsssetStoreError:
		return returnErrorWithCode(500, err, message, logCtx)

	case *shared.UserInputError:
		return returnErrorWithCode(400, err, message, logCtx)

	case *entitystore.EntityStoreError:
		return returnErrorWithCode(500, err, message, logCtx)

	default:
		return returnErrorWithCode(500, err, message, logCtx)
	}
}

func returnErrorWithCode(statusCode int, err error, message string, logCtx zerolog.Context) (events.APIGatewayProxyResponse, error) {
	shared.GetLogger(logCtx).Err(err).Msg(message)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}, nil
}
