package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	talktomeartcreate "talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/talktome"
	"talktome.com/internal/user"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// ENV VAR init
	openAIToken := shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	conversationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
	userTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")
	conversationClipBucket := shared.MustReadEnvVar("TALKTOME_CONVERSATION_CLIP_BUCKET")

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
	textGen := conversation.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewAWSPollySpeechGenerator(pollyClient)
	convStorage := conversation.NewAWSStorageCtx(dynamoDBClient, conversationDynamoDBTable, s3, conversationClipBucket)
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userTable)

	ctx := talktome.NewContext(textGen, speechGen, convStorage, userStorage)

	log.Info().Msg("starting 'start art conversation' lambda")
	lambda.Start(talktomeartcreate.HandlerCtx{Ctx: ctx}.AWSHandler)
}
