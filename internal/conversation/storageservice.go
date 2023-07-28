package conversation

import (
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageService interface {
	FindConversation(uuid string) (*Conversation, error)
	StoreConversation(conversation Conversation) error
	StoreClip(clip *os.File) error
}

type AWSStorageCtx struct {
	table          string
	dynamoDBClient *dynamodb.DynamoDB

	s3Client   *s3.S3
	bucketName string
}

func NewAWSStorageCtx(dynamoDB *dynamodb.DynamoDB, table string, s3Client *s3.S3, bucketName string) *AWSStorageCtx {
	return &AWSStorageCtx{
		table:          table,
		dynamoDBClient: dynamoDB,
		s3Client:       s3Client,
		bucketName:     bucketName,
	}
}
