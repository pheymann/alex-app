package assetstore

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx AWSS3Context) GenerateTemporaryAccessURL(assetUUID string, logCtx zerolog.Context) (string, error) {
	shared.GetLogger(logCtx).Debug().Msg("generate temporary access url")

	input := &s3.GetObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(assetUUID),
	}

	audioClipRequest, _ := ctx.s3Client.GetObjectRequest(input)

	urlValidFor := 30 * time.Minute
	url, err := audioClipRequest.Presign(urlValidFor)
	if err != nil {
		return "", &AsssetStoreError{err, fmt.Sprintf("failed to generate access url for asset %s", assetUUID)}
	}

	shared.GetLogger(logCtx).Debug().Msgf("generated access url: %s", url)
	return url, nil
}
