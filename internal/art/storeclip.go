package art

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ctx *StorageCtx) StoreClip(clipUUID string, audioContent []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(clipUUID + ".wav"),
		Body:   bytes.NewReader(audioContent),
	}

	if _, err := ctx.s3Client.PutObject(input); err != nil {
		return fmt.Errorf("failed to store WAV file for clip %s: %w", clipUUID, err)
	}

	return nil
}
