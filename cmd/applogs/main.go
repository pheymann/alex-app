package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/applogs"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("starting 'app logs' lambda")
	lambda.Start(applogs.AWSHandler)
}
