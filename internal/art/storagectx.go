package art

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageCtx struct {
	table          string
	dynamoDBClient *dynamodb.DynamoDB

	s3Client   *s3.S3
	bucketName string
}

func NewStorageCtx(dynamoDB *dynamodb.DynamoDB, table string, s3Client *s3.S3, bucketName string) *StorageCtx {
	return &StorageCtx{
		table:          table,
		dynamoDBClient: dynamoDB,
		s3Client:       s3Client,
		bucketName:     bucketName,
	}
}
