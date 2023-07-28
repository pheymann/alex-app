package user

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type StorageService interface {
	StoreUser(user User) error
	FindUser(uuid string) (*User, error)
}

type AWSStorageCtx struct {
	table          string
	dynamoDBClient *dynamodb.DynamoDB
}

func NewAWSStorageCtx(dynamoDB *dynamodb.DynamoDB, table string) *AWSStorageCtx {
	return &AWSStorageCtx{
		table:          table,
		dynamoDBClient: dynamoDB,
	}
}
