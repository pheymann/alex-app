package art

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (ctx *StorageCtx) StoreArtPresentation(presentation ArtPresentation) error {
	item, err := dynamodbattribute.MarshalMap(presentation)
	if err != nil {
		return fmt.Errorf("failed to marshal presentation :%w", err)
	}

	_, err = ctx.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ArtPiecePresentation"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to store presentation :%w", err)
	}

	return nil
}
