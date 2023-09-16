package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/zerolog"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/shared"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	operation := flag.String("operation", "", "--operation <operation>")

	userUUID := flag.String("user-uuid", "", "--user-uuid <uuid>")

	convUUID := flag.String("conv-uuid", "", "--conv-uuid <uuid>")
	message := flag.String("message", "", "--message <message>")

	artContext := flag.String("art-context", "", "--art-context <full name>")

	flag.Parse()

	if *operation == "" {
		panic("missing operation")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	switch *operation {
	case "create-art":
		createArtConversation(*userUUID, artContext, sess)
		return
	case "continue":
		continueConversation(*userUUID, convUUID, message, sess)
		return
	case "list-all":
		listAllConversations(*userUUID, sess)
		return
	case "get":
		getConversation(*userUUID, convUUID, sess)
		return

	default:
		panic(fmt.Sprintf("unknown operation: %s", *operation))
	}
}

func createArtConversation(userUUID string, artContext *string, sess *session.Session) {
	if *artContext == "" {
		panic("missing artist context")
	}

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Body:           fmt.Sprintf(`{"artContext": "%s"}`, *artContext),
		RequestContext: shared.NewAwsTestRequestContext(userUUID),
	}

	response, err := startartconversation.
		UnsafeNewHandlerCtx(
			sess,
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE"),
			shared.MustReadEnvVar("TALKTOME_USER_TABLE"),
			shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN"),
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_CLIP_BUCKET"),
		).
		AWSHandler(context.TODO(), event)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", response)
}

func continueConversation(userUUID string, convUUID, message *string, sess *session.Session) {
	if *message == "" {
		panic("if 'conv-uuid' is set you have to provide a message")
	} else if *convUUID == "" {
		panic("missing conversation uuid")
	}

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     "POST",
		Body:           fmt.Sprintf(`{"question": "%s"}`, *message),
		RequestContext: shared.NewAwsTestRequestContext(userUUID),
		PathParameters: map[string]string{
			"uuid": *convUUID,
		},
	}

	response, err := continueconversation.
		UnsafeNewHandlerCtx(
			sess,
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE"),
			shared.MustReadEnvVar("TALKTOME_USER_TABLE"),
			shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN"),
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_CLIP_BUCKET"),
		).
		AWSHandler(context.TODO(), event)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", response)
}

func listAllConversations(userUUID string, sess *session.Session) {
	event := events.APIGatewayProxyRequest{
		HTTPMethod:     "GET",
		RequestContext: shared.NewAwsTestRequestContext(userUUID),
	}

	response, err := listconversations.
		UnsafeNewHandlerCtx(
			sess,
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE"),
			shared.MustReadEnvVar("TALKTOME_USER_TABLE"),
		).
		AWSHandler(context.TODO(), event)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", response)
}

func getConversation(userUUID string, convUUID *string, sess *session.Session) {
	if *convUUID == "" {
		panic("missing conversation uuid")
	}

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     "GET",
		RequestContext: shared.NewAwsTestRequestContext(userUUID),
		PathParameters: map[string]string{
			"uuid": *convUUID,
		},
	}

	response, err := getconversation.
		UnsafeNewHandlerCtx(
			sess,
			shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE"),
			shared.MustReadEnvVar("TALKTOME_USER_TABLE"),
		).
		AWSHandler(context.TODO(), event)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", response)
}
