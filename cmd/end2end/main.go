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
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/talktome"
	"talktome.com/internal/user"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	operation := flag.String("operation", "", "--operation <operation>")

	userUUID := flag.String("user-uuid", "", "--user-uuid <uuid>")

	convUUID := flag.String("conv-uuid", "", "--conv-uuid <uuid>")
	message := flag.String("message", "", "--message <message>")

	artContext := flag.String("art-context", "", "--art-context <full name>")

	flag.Parse()

	if *operation == "" {
		panic("missing operation")
	}

	switch *operation {
	case "create-art":
		createArtConversation(*userUUID, artContext)
		return
	case "continue":
		continueConversation(*userUUID, convUUID, message)
		return
	case "list-all":
		listAllConversations(*userUUID)
		return
	case "get":
		getConversation(*userUUID, convUUID)
		return

	default:
		panic(fmt.Sprintf("unknown operation: %s", *operation))
	}
}

func createArtConversation(userUUID string, artContext *string) {
	if *artContext == "" {
		panic("missing artist context")
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

	log.Info().Msg("creating new art conversation")
	conv, err := startartconversation.Handle(userUUID, *artContext, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", *conv)
}

func continueConversation(userUUID string, convUUID, message *string) {
	if *message == "" {
		panic("if 'conv-uuid' is set you have to provide a message")
	} else if *convUUID == "" {
		panic("missing conversation uuid")
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

	log.Info().Msg("continue conversation")
	conv, err := continueconversation.Handle(userUUID, *convUUID, *message, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", *conv)
}

func listAllConversations(userUUID string) {
	// ENV VAR init
	conversationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
	userDynamoDBTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")

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
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userDynamoDBTable)

	log.Info().Msg("list all conversations")
	conv, err := listconversations.Handle(userUUID, userStorage, convStorage)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", conv)
}

func getConversation(userUUID string, convUUID *string) {
	if *convUUID == "" {
		panic("missing conversation uuid")
	}

	// ENV VAR init
	conversationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
	userDynamoDBTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")

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
	userStorage := user.NewAWSStorageCtx(dynamoDBClient, userDynamoDBTable)

	log.Info().Msg("get conversation")
	conv, err := getconversation.Handle(userUUID, *convUUID, userStorage, convStorage)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", conv)
}
