package entitystore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx AWSDynamoDBCtx[E]) Save(entity E, logCtx zerolog.Context) error {
	shared.GetLogger(logCtx).Debug().Msg("save entity")

	item, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		return NewEntityStoreError("failed to marshal entity", err)
	}

	_, err = ctx.dynamoDBClient.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{
			TransactItems: []*dynamodb.TransactWriteItem{
				{
					Put: &dynamodb.Put{
						TableName: aws.String(ctx.table),
						Item:      item,
					},
				},
			},
		},
	)
	if err != nil {
		return NewEntityStoreError("failed to store conversation", err)
	}

	return nil
}
