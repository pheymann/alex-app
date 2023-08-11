package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
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

	dynamoDBClient := dynamodb.New(sess)

	// internal init
	convStorage := conversation.NewAWSStorageCtx(dynamoDBClient, conversationDynamoDBTable, nil, "")
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userTable)

	log.Info().Str("conversation_table", conversationDynamoDBTable).Str("user_table", userTable).Msg("starting 'list conversations' lambda")
	lambda.Start(listconversations.HandlerCtx{UserStorage: userStorage, ConvStorage: convStorage}.AWSHandler)
}
