package assetstore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx AWSS3Context) Save(file *os.File, logCtx zerolog.Context) (string, error) {
	shared.GetLogger(logCtx).Debug().Msg("save file")

	name := file.Name()
	key := filepath.Base(name)

	input := &s3.PutObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    &key,
		Body:   file,
	}

	if _, err := ctx.s3Client.PutObject(input); err != nil {
		return key, &AsssetStoreError{err, fmt.Sprintf("failed to save file %s", name)}
	}

	shared.GetLogger(logCtx).Debug().Msgf("saved file %s under key: %s", name, key)

	return key, nil
}
