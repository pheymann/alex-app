package art

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (ctx *StorageCtx) StoreArtConversation(conversation ArtConversation) error {
	item, err := dynamodbattribute.MarshalMap(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation :%w", err)
	}

	_, err = ctx.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ctx.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to store conversation :%w", err)
	}

	return nil
}
