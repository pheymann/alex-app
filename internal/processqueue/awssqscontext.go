package processqueue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AWSSQSContext struct {
	queueURL  string
	sqsClient *sqs.SQS
}

func NewAWSSQSContext(sqsClient *sqs.SQS, queueURL string) *AWSSQSContext {
	return &AWSSQSContext{
		queueURL:  queueURL,
		sqsClient: sqsClient,
	}
}
