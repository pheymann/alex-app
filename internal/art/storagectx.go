package art

import "github.com/aws/aws-sdk-go/service/dynamodb"

type StorageCtx struct {
	table          string
	dynamoDBClient *dynamodb.DynamoDB
}

func NewStorageCtx(dynamoDB *dynamodb.DynamoDB, table string) *StorageCtx {
	return &StorageCtx{
		table:          table,
		dynamoDBClient: dynamoDB,
	}
}
