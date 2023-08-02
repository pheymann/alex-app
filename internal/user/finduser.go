package user

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func (ctx *AWSStorageCtx) FindUser(uuid string) (*User, error) {
	log.Debug().Str("uuid", uuid).Msg("find user")

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

	var user = User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user :%w", err)
	}

	return &user, nil
}
