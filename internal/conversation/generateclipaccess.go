package conversation

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ctx *AWSStorageCtx) GenerateClipAccess(audioClipUUID string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(audioClipUUID),
	}

	audioClipRequest, _ := ctx.s3Client.GetObjectRequest(input)

	// TODO: configure expiration time
	url, err := audioClipRequest.Presign(3 * time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate access url for audio clip %s: %w", audioClipUUID, err)
	}

	return url, nil
}
