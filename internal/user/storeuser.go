package user

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func (ctx *AWSStorageCtx) StoreUser(user User) error {
	log.Debug().Str("user_uuid", user.ID).Msg("store user")

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user :%w", err)
	}

	_, err = ctx.dynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ctx.table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to store user :%w", err)
	}

	return nil
}
