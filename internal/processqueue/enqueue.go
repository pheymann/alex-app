package processqueue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog"
	"talktome.com/internal/shared"
)

func (ctx *AWSSQSContext) Enqueue(task Task, logCtx zerolog.Context) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		return NewProcessQueueError("Failed to marshal task", err)
	}

	taskStr := string(taskJson)

	input := sqs.SendMessageInput{
		MessageBody: &taskStr,
		QueueUrl:    &ctx.queueURL,
	}

	if _, err := ctx.sqsClient.SendMessage(&input); err != nil {
		return NewProcessQueueError("Failed to enqueue task", err)
	}
	shared.GetLogger(logCtx).Debug().Msg("enqueued task")

	return nil
}
