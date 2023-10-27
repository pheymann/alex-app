package assistant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/assetstore"
	"talktome.com/internal/conversation"
	"talktome.com/internal/entitystore"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

type HandlerCtx struct {
	ConversationStore entitystore.EntityStore[conversation.Conversation]
	UserStore         entitystore.EntityStore[user.User]
	AudioClipStore    assetstore.AssetStore
	TextGen           textgeneration.TextGenerationService
	SpeechGen         speechgeneration.SpeechGenerator
}

func (handlerCtx HandlerCtx) AWSHandler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		logCtx := log.With().Str("lambda", "assistant")

		var task processqueue.Task
		if err := json.Unmarshal([]byte(record.Body), &task); err != nil {
			return &shared.UserInputError{Message: fmt.Sprintf("couldn't parse body: %e", err)}
		}

		logCtx = logCtx.Str("user_uuid", task.UserUUID)
		logCtx = logCtx.Str("conversation_uuid", task.ConversationUUID)
		logCtx = logCtx.Str("language", string(task.Language))

		convCtx := conversation.Context{
			ConversationUUID:  task.ConversationUUID,
			UserUUID:          task.UserUUID,
			Language:          task.Language,
			LogCtx:            logCtx,
			ConversationStore: handlerCtx.ConversationStore,
			UserStore:         handlerCtx.UserStore,
			AudioClipStore:    handlerCtx.AudioClipStore,
			TextGen:           handlerCtx.TextGen,
			SpeechGen:         handlerCtx.SpeechGen,
		}

		if err := Handle(convCtx, task.Message); err != nil {
			return &shared.InternalError{
				Cause:   err,
				Message: "failed to handle task",
			}
		}
	}

	return nil
}

func UnsafeNewHandlerCtx(
	sess *session.Session,
	conversationDynamoDBTable string,
	userTable string,
	openAIToken string,
	conversationClipBucket string,
) HandlerCtx {
	dynamoDBClient := dynamodb.New(sess)
	s3 := s3.New(sess)
	pollyClient := polly.New(sess)

	// internal init
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewAWSPollySpeechGenerator(pollyClient)
	convStorage := entitystore.NewAWSDynamoDBCtx[conversation.Conversation](dynamoDBClient, conversationDynamoDBTable)
	userStorage := entitystore.NewAWSDynamoDBCtx[user.User](dynamoDBClient, userTable)
	audioClipStore := assetstore.NewAWSS3Context(s3, conversationClipBucket)

	return HandlerCtx{
		ConversationStore: convStorage,
		UserStore:         userStorage,
		AudioClipStore:    audioClipStore,
		TextGen:           textGen,
		SpeechGen:         speechGen,
	}
}
