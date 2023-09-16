package assetstore

import "github.com/aws/aws-sdk-go/service/s3"

type AWSS3Context struct {
	s3Client   *s3.S3
	bucketName string
}

func NewAWSS3Context(s3Client *s3.S3, bucketName string) *AWSS3Context {
	return &AWSS3Context{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}
