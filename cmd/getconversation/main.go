package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/shared"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	ssmClient := ssm.New(sess)

	ctx := getconversation.UnsafeNewHandlerCtx(
		sess,
		shared.MustReadParameter("talktome-table-conversation", ssmClient),
		shared.MustReadParameter("talktome-table-user", ssmClient),
		shared.MustReadParameter("talktome-bucket-conversation-clips", ssmClient),
	)

	log.Info().Msg("starting 'get conversation' lambda")
	lambda.Start(ctx.AWSHandler)
}
