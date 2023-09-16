package entitystore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx *AWSDynamoDBCtx[E]) Find(uuid string, logCtx zerolog.Context) (*E, error) {
	shared.GetLogger(logCtx).Debug().Msg("find entity")

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
		return nil, NewEntityStoreError("failed to load entity", err)
	}

	if result.Item == nil || len(result.Item) == 0 {
		return nil, nil
	}

	var entity E
	err = dynamodbattribute.UnmarshalMap(result.Item, &entity)
	if err != nil {
		return nil, NewEntityStoreError("failed to unmarshal entity", err)
	}

	return &entity, nil
}
