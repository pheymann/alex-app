package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/user"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// ENV VAR init
	conversationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
	userTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	dynamoDBClient := dynamodb.New(sess)

	// internal init
	convStorage := conversation.NewAWSStorageCtx(dynamoDBClient, conversationDynamoDBTable, nil, "")
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userTable)

	log.Info().Msg("starting 'get conversation' lambda")
	lambda.Start(getconversation.HandlerCtx{UserStorage: userStorage, ConvStorage: convStorage}.AWSHandler)
}
