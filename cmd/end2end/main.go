package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/talktomeartcreate"
	"talktome.com/internal/cmd/talktomecontinue"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/talktome"
	"talktome.com/internal/user"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	userUUID := flag.String("user-uuid", "", "--user-uuid <uuid>")

	convUUID := flag.String("conv-uuid", "", "--conv-uuid <uuid>")
	message := flag.String("message", "", "--message <message>")

	artistName := flag.String("artist", "", "--artist <full name>")
	artPiece := flag.String("art-piece", "", "--art-piece <full name>")

	flag.Parse()

	if *userUUID == "" {
		panic("missing user uuid")
	}

	if *convUUID == "" {
		if *artistName == "" {
			panic("missing artist name")
		} else if *artPiece == "" {
			panic("missing art piece name")
		}
	} else {
		if *message == "" {
			panic("if 'conv-uuid' is set you have to provide a message")
		}
	}

	// ENV VAR init
	openAIToken := shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	conversationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
	userDynamoDBTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")
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
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userDynamoDBTable)

	ctx := talktome.NewContext(textGen, speechGen, convStorage, userStorage)

	if *convUUID == "" {
		log.Info().Msg("creating new conversation")
		conv, err := talktomeartcreate.Handle(*userUUID, *artistName, *artPiece, ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v", *conv)
	} else {
		message, err := talktomecontinue.Handle(*userUUID, *convUUID, *message, ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v", *message)
	}
}
