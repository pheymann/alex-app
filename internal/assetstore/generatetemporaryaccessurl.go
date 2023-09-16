package assetstore

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx AWSS3Context) GenerateTemporaryAccessURL(assetUUID string, logCtx zerolog.Context) (string, *time.Time, error) {
	shared.GetLogger(logCtx).Debug().Msg("generate temporary access url")

	input := &s3.GetObjectInput{
		Bucket: aws.String(ctx.bucketName),
		Key:    aws.String(assetUUID),
	}

	audioClipRequest, _ := ctx.s3Client.GetObjectRequest(input)

	// TODO: configure expiration time
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return "", nil, &AsssetStoreError{err, "failed to load UTC location"}
	}

	urlValidFor := 3 * 24 * time.Hour
	expirationDate := time.Now().In(location).Add(urlValidFor - 1*time.Hour)
	url, err := audioClipRequest.Presign(urlValidFor)
	if err != nil {
		return "", nil, &AsssetStoreError{err, fmt.Sprintf("failed to generate access url for asset %s", assetUUID)}
	}

	shared.GetLogger(logCtx).Debug().Msgf("generated access url: %s", url)
	return url, &expirationDate, nil
}
