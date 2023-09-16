package entitystore

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type AWSDynamoDBCtx[E any] struct {
	table          string
	dynamoDBClient *dynamodb.DynamoDB
}

func NewAWSDynamoDBCtx[E any](dynamoDB *dynamodb.DynamoDB, table string) *AWSDynamoDBCtx[E] {
	return &AWSDynamoDBCtx[E]{
		table:          table,
		dynamoDBClient: dynamoDB,
	}
}
