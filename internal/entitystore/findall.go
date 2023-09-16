package entitystore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx *AWSDynamoDBCtx[E]) FindAll(uuids []string, logCtx zerolog.Context) ([]E, error) {
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
	logCtx.Array("uuids", logUUIDArray)
	shared.GetLogger(logCtx).Debug().Msg("find all entities")

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			ctx.table: {
				Keys: keys,
			},
		},
	}

	result, err := ctx.dynamoDBClient.BatchGetItem(input)
	if err != nil {
		return nil, NewEntityStoreError("failed to load enitity", err)
	}

	// TODO: how to handle unprocessed keys?
	entities := make([]E, len(result.Responses[ctx.table]))
	err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[ctx.table], &entities)
	if err != nil {
		return nil, NewEntityStoreError("failed to unmarshal entity", err)
	}

	shared.GetLogger(logCtx).Debug().Msgf("found %d entities", len(entities))

	return entities, nil
}
