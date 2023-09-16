package applogs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/shared"
)

type appLogs struct {
	LogEntries []appLogEntry `json:"logEntries"`
}

type appLogEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func AWSHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userUUID, error := shared.ExtractUserUUID(event)
	if error != nil {
		log.Err(error).Msg("failed to extract user uuid; continue without it")
	}

	var logs appLogs
	if err := json.Unmarshal([]byte(event.Body), &logs); err != nil {
		log.Err(err).Msg("couldn't parse body and abort logging")
		return staticSuccessResponse, nil
	}

	logCtx := log.With().Str("user_uuid", userUUID).Str("system", "app")

	for _, logEntry := range logs.LogEntries {
		logger := logCtx.Str("app_timestamp", logEntry.Timestamp).Logger()

		switch logEntry.Level {
		case "debug":
			logger.Debug().Msg(logEntry.Message)
		case "info":
			logger.Info().Msg(logEntry.Message)
		case "warn":
			logger.Warn().Msg(logEntry.Message)
		case "error":
			logger.Error().Msg(logEntry.Message)
		case "fatal":
			logger.Fatal().Msg(logEntry.Message)
		default:
			logger.Info().Msg(logEntry.Message)
		}
	}

	return staticSuccessResponse, nil
}

var (
	staticSuccessResponse = events.APIGatewayProxyResponse{
		StatusCode: 201,
	}
)
