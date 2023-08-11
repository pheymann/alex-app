package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/talktome"
	"talktome.com/internal/user"
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

	conversationDynamoDBTable := shared.MustReadParameter("talktome-table-conversation", ssmClient)
	userTable := shared.MustReadParameter("talktome-table-user", ssmClient)
	openAIToken := shared.MustReadParameter("talktome-secret-openai-token", ssmClient)
	conversationClipBucket := shared.MustReadParameter("talktome-bucket-conversation-clips", ssmClient)

	dynamoDBClient := dynamodb.New(sess)
	s3 := s3.New(sess)
	pollyClient := polly.New(sess)

	// internal init
	textGen := conversation.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewAWSPollySpeechGenerator(pollyClient)
	convStorage := conversation.NewAWSStorageCtx(dynamoDBClient, conversationDynamoDBTable, s3, conversationClipBucket)
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userTable)

	ctx := talktome.NewContext(textGen, speechGen, convStorage, userStorage)

	log.Info().
		Str("conversation_table", conversationDynamoDBTable).
		Str("user_table", userTable).
		Str("clip_bucket", conversationClipBucket).
		Msg("starting 'continue conversation' lambda")
	lambda.Start(continueconversation.HandlerCtx{Ctx: ctx}.AWSHandler)
}
