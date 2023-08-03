package conversation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (ctx *AWSStorageCtx) FindAllConversations(uuids []string) ([]Conversation, error) {
	logUUIDArray := zerolog.Arr()
	keys := make([]map[string]*dynamodb.AttributeValue, len(uuids))

	for index, uuid := range uuids {
		logUUIDArray.Str(uuid)

		keys[index] = map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(uuid),
			},
		}
	}

	log.Debug().Array("conversation_uuids", logUUIDArray).Msg("find all conversation")

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			ctx.table: {
				Keys: keys,
			},
		},
	}

	result, err := ctx.dynamoDBClient.BatchGetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversations :%w", err)
	}

	// TODO: how to handle unprocessed keys?
	conversations := make([]Conversation, len(result.Responses[ctx.table]))
	err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[ctx.table], &conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversations :%w", err)
	}

	return conversations, nil
}
