package conversation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (ctx *AWSStorageCtx) StoreClip(clipFile *os.File) (string, error) {
	key := filepath.Base(clipFile.Name())

	input := &s3.PutObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    &key,
		Body:   clipFile,
	}

	if _, err := ctx.s3Client.PutObject(input); err != nil {
		return key, fmt.Errorf("failed to store WAV file for clip %s: %w", clipFile.Name(), err)
	}

	return key, nil
}
