package conversation

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ctx *AWSStorageCtx) GenerateClipAccess(audioClipUUID string) (string, *time.Time, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(audioClipUUID),
	}

	audioClipRequest, _ := ctx.s3Client.GetObjectRequest(input)

	// TODO: configure expiration time
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return "", nil, fmt.Errorf("failed to load UTC location: %w", err)
	}

	urlValidFor := 3 * 24 * time.Hour
	expirationDate := time.Now().In(location).Add(urlValidFor - 1*time.Hour)
	url, err := audioClipRequest.Presign(urlValidFor)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate access url for audio clip %s: %w", audioClipUUID, err)
	}

	return url, &expirationDate, nil
}
