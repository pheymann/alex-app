package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"talktome.com/internal/art"
	"talktome.com/internal/cmd/talktome"
	"talktome.com/internal/shared"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
)

func main() {
	artistName := flag.String("artist", "", "--artist <full name>")
	artPieceName := flag.String("art-piece", "", "--art-piece <full name>")

	flag.Parse()

	if *artistName == "" {
		panic("missing artist name")
	}

	if *artPieceName == "" {
		panic("missing art piece name")
	}

	artPiece := art.ArtPiece{
		ArtistName: *artistName,
		Name:       *artPieceName,
	}

	// ENV VAR init
	openAIToken := shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
	artPresentationDynamoDBTable := shared.MustReadEnvVar("TALKTOME_ART_PRESENTATION_TABLE")

	// AWS init
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		panic(err)
	}

	dynamoDBClient := dynamodb.New(sess)
	s3 := s3.New(sess)
	pollyClient := polly.New(sess)

	// internal init
	textGen := textgeneration.NewOpenAIGenerator(openAIToken)
	speechGen := speechgeneration.NewPollySpeechGenerator(pollyClient)
	artStorage := art.NewStorageCtx(dynamoDBClient, artPresentationDynamoDBTable, s3, "talktome-artaudioclips")

	talktome := talktome.NewTalkToMe(textGen, speechGen, artStorage)

	conversation, _, err := talktome.TalkToMeArt(artPiece, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("=> %v", *conversation)

	var prompt = "Can you show me 2 or 3 interesting aspects of that painting?"

	// fmt.Print("Do you have a question? ")
	// fmt.Scanln(prompt)

	convContinued, clip, err := talktome.TalkToMeArt(artPiece, &prompt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("=> %v", *convContinued)

	file, err := os.Create("prompt.mp3")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the byte array to the file
	_, err = file.Write(clip)
	if err != nil {
		panic(err)
	}
}
