package art

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ctx *AWSStorageCtx) StoreClip(clipFile *os.File) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(clipFile.Name()),
		Body:   clipFile,
	}

	if _, err := ctx.s3Client.PutObject(input); err != nil {
		return fmt.Errorf("failed to store WAV file for clip %s: %w", clipFile.Name(), err)
	}

	return nil
}
