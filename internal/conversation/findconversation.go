package conversation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func (ctx *AWSStorageCtx) FindConversation(uuid string) (*Conversation, error) {
	log.Debug().Str("conversation_uuid", uuid).Msg("find conversation")

	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(uuid),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(ctx.table),
		Key:       key,
	}

	result, err := ctx.dynamoDBClient.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load presentation :%w", err)
	}

	if result.Item == nil || len(result.Item) == 0 {
		return nil, nil
	}

	var conversation = Conversation{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation :%w", err)
	}

	return &conversation, nil
}