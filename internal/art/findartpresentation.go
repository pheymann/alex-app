package art

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (ctx *StorageCtx) FindArtPresentation(piece ArtPiece) (*ArtPresentation, error) {
	presentation := ArtPresentation{Piece: piece}
	key, err := dynamodbattribute.MarshalMap(presentation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal presentation :%w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(ctx.table),
		Key:       key,
	}

	result, err := ctx.dynamoDBClient.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load presentation :%w", err)
	}

	if result.Item == nil || len(result.Item) == 0 {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &presentation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal presentation :%w", err)
	}

	return &presentation, nil
}
