package main

import (
	"bytes"
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/applogs"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/testutil"
	"talktome.com/internal/user"
)

const (
	longText = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.

	Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.

	Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis.

	At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, At accusam aliquyam diam diam dolore dolores duo eirmod eos erat, et nonumy sed tempor et et invidunt justo labore Stet clita ea et gubergren, kasd magna no rebum. sanctus sea sed takimata ut vero voluptua. est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur`
)

var (
	mockTextGen = &testutil.MockTextGeneration{
		GeneratedMessage: longText,
		Timeout:          5 * time.Second,
	}

	mockSpeechGen = &testutil.MockSpeechGeneration{
		TestFile: "assets/prompt.mp3",
		Timeout:  3 * time.Second,
	}

	mockConversationStore = testutil.MockConversationStore(map[string]*conversation.Conversation{})

	testUserRequestCtx = shared.NewAwsTestRequestContext("1")

	mockUserStorage = testutil.MockUserStore(map[string]*user.User{
		"1": {
			ID:                "1",
			ConversationUUIDs: []string{},
		},
	})

	mockAudioClipStore = &testutil.MockAssetStore{
		ClipKey:      "prompt.mp3",
		PresignedUrl: "/aws/presigned/prompt.mp3",
	}
)

func handleStartArtConversation(ctx startartconversation.HandlerCtx) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		headers := make(map[string]string)
		for k, v := range r.Header {
			headers[k] = v[0]
		}

		event := events.APIGatewayProxyRequest{
			HTTPMethod:     r.Method,
			Body:           buf.String(),
			Headers:        headers,
			RequestContext: testUserRequestCtx,
		}

		response, err := ctx.AWSHandler(context.TODO(), event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		if _, err := w.Write([]byte(response.Body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleContinueConversation(ctx continueconversation.HandlerCtx) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		headers := make(map[string]string)
		for k, v := range r.Header {
			headers[k] = v[0]
		}

		event := events.APIGatewayProxyRequest{
			HTTPMethod: r.Method,
			PathParameters: map[string]string{
				"uuid": vars["id"],
			},
			Body:           buf.String(),
			Headers:        headers,
			RequestContext: testUserRequestCtx,
		}

		response, err := ctx.AWSHandler(context.TODO(), event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		if _, err := w.Write([]byte(response.Body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGetConversation(ctx getconversation.HandlerCtx) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		event := events.APIGatewayProxyRequest{
			HTTPMethod: r.Method,
			PathParameters: map[string]string{
				"uuid": vars["id"],
			},
			RequestContext: testUserRequestCtx,
		}

		response, err := ctx.AWSHandler(context.TODO(), event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		if _, err := w.Write([]byte(response.Body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleListConversations(ctx listconversations.HandlerCtx) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		event := events.APIGatewayProxyRequest{
			HTTPMethod:     r.Method,
			RequestContext: testUserRequestCtx,
		}

		response, err := ctx.AWSHandler(context.TODO(), event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		if _, err := w.Write([]byte(response.Body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleAppLogs(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     r.Method,
		Body:           buf.String(),
		RequestContext: testUserRequestCtx,
	}

	response, err := applogs.AWSHandler(context.Background(), event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "assets/prompt.mp3")
	}
}

func main() {
	mode := flag.String("mode", "", "--mode prod|mock")

	flag.Parse()

	if mode == nil {
		panic("mode is not defined")
	}

	var startArtConvCtx startartconversation.HandlerCtx
	var continueConvCtx continueconversation.HandlerCtx
	var listConvCtx listconversations.HandlerCtx
	var getConvCtx getconversation.HandlerCtx

	if *mode == "mock" {
		startArtConvCtx = startartconversation.HandlerCtx{
			ConversationStore: mockConversationStore,
			UserStore:         mockUserStorage,
			AudioClipStore:    mockAudioClipStore,
			TextGen:           mockTextGen,
			SpeechGen:         mockSpeechGen,
		}

		continueConvCtx = continueconversation.HandlerCtx{
			ConversationStore: mockConversationStore,
			UserStore:         mockUserStorage,
			AudioClipStore:    mockAudioClipStore,
			TextGen:           mockTextGen,
			SpeechGen:         mockSpeechGen,
		}

		listConvCtx = listconversations.HandlerCtx{
			ConversationStore: mockConversationStore,
			UserStore:         mockUserStorage,
		}

		getConvCtx = getconversation.HandlerCtx{
			ConversationStore: mockConversationStore,
			UserStore:         mockUserStorage,
		}
	} else if *mode == "prod" {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-central-1"),
		})
		if err != nil {
			panic(err)
		}

		conversationTable := shared.MustReadEnvVar("TALKTOME_CONVERSATION_TABLE")
		userTable := shared.MustReadEnvVar("TALKTOME_USER_TABLE")
		openAIToken := shared.MustReadEnvVar("TALKTOME_OPEN_AI_TOKEN")
		clipBucket := shared.MustReadEnvVar("TALKTOME_CONVERSATION_CLIP_BUCKET")

		startArtConvCtx = startartconversation.UnsafeNewHandlerCtx(
			sess,
			conversationTable,
			userTable,
			openAIToken,
			clipBucket,
		)

		continueConvCtx = continueconversation.UnsafeNewHandlerCtx(
			sess,
			conversationTable,
			userTable,
			openAIToken,
			clipBucket,
		)

		listConvCtx = listconversations.UnsafeNewHandlerCtx(
			sess,
			conversationTable,
			userTable,
		)

		getConvCtx = getconversation.UnsafeNewHandlerCtx(
			sess,
			conversationTable,
			userTable,
		)
	} else {
		panic("unknown mode: " + *mode)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	router := mux.NewRouter()

	router.HandleFunc("/api/conversation/create/art", handleStartArtConversation(startArtConvCtx)).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/list", handleListConversations(listConvCtx)).Methods(http.MethodGet)
	router.HandleFunc("/api/conversation/{id}/continue", handleContinueConversation(continueConvCtx)).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/{id}", handleGetConversation(getConvCtx)).Methods(http.MethodGet)
	router.HandleFunc("/api/app/logs", handleAppLogs).Methods(http.MethodPost)
	router.HandleFunc("/aws/presigned/{id}", fileHandler).Methods(http.MethodGet)

	port := ":8080"
	log.Info().Msgf("Server running on port %s and targeting %s environment", port, *mode)
	log.Fatal().Err(http.ListenAndServe(port, router))
}
