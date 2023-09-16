package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"talktome.com/internal/cmd/applogs"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

// Mocks
type mockTextGeneration struct {
	generatedMessage string
}

func (generator *mockTextGeneration) GenerateNextMessage(
	messageHistory []textgeneration.BasicMessage,
	logCtx zerolog.Context,
) (*textgeneration.BasicMessage, error) {
	// simulate OpenAI generation time
	time.Sleep(5 * time.Second)

	return &textgeneration.BasicMessage{
		Role: textgeneration.RoleAssistent,
		Text: generator.generatedMessage,
	}, nil
}

type mockSpeechGeneration struct{}

func (generator *mockSpeechGeneration) GenerateSpeechClip(
	title string,
	text string,
	logCtx zerolog.Context,
) (*os.File, error) {
	file, err := os.Open("assets/prompt.mp3")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	copyFile, err := os.CreateTemp("", "speechclip*.mp3")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(copyFile, file); err != nil {
		return nil, err
	}

	// simulate Polly generation time
	time.Sleep(3 * time.Second)

	return copyFile, nil
}

type mockEntityStore[E any] struct {
	localStore   map[string]*E
	makeDeepCopy func(entity *E) *E
	getID        func(entity E) string
}

func (ctx *mockEntityStore[E]) Find(uuid string, logCtx zerolog.Context) (*E, error) {
	if entity, ok := ctx.localStore[uuid]; ok {
		return ctx.makeDeepCopy(entity), nil
	}
	return nil, nil
}

func (ctx *mockEntityStore[E]) FindAll(uuids []string, logCtx zerolog.Context) ([]E, error) {
	var entities []E
	for _, uuid := range uuids {
		if entity, ok := ctx.localStore[uuid]; ok {
			entities = append(entities, *entity)
		}
	}

	return entities, nil
}

func (ctx *mockEntityStore[E]) Save(entity E, logCtx zerolog.Context) error {
	ctx.localStore[ctx.getID(entity)] = ctx.makeDeepCopy(&entity)
	return nil
}

type mockAssetStore struct{}

func (ctx *mockAssetStore) Save(file *os.File, logCtx zerolog.Context) (string, error) {
	return "prompt.mp3", nil
}

func (ctx *mockAssetStore) GenerateTemporaryAccessURL(audioClipUUID string, logCtx zerolog.Context) (string, *time.Time, error) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return "", nil, err
	}

	urlValidFor := 1 * time.Minute
	expirationDate := time.Now().In(location).Add(urlValidFor)

	return "/aws/presigned/prompt.mp3", &expirationDate, nil
}

const (
	longText = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.

	Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.

	Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis.

	At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, At accusam aliquyam diam diam dolore dolores duo eirmod eos erat, et nonumy sed tempor et et invidunt justo labore Stet clita ea et gubergren, kasd magna no rebum. sanctus sea sed takimata ut vero voluptua. est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur`
)

var (
	mockTextGen = &mockTextGeneration{
		generatedMessage: longText,
	}

	mockSpeechGen = &mockSpeechGeneration{}

	mockConversationStore = &mockEntityStore[conversation.Conversation]{
		localStore: map[string]*conversation.Conversation{},
		makeDeepCopy: func(conv *conversation.Conversation) *conversation.Conversation {
			convCopy := conversation.Conversation{
				ID:       conv.ID,
				Metadata: conv.Metadata,
				Messages: []conversation.Message{},
			}

			// deep copy
			convCopy.Messages = append(convCopy.Messages, conv.Messages...)

			return &convCopy
		},
		getID: func(conv conversation.Conversation) string {
			return conv.ID
		},
	}

	testUserRequestCtx = shared.NewAwsTestRequestContext("1")

	mockUserStorage = &mockEntityStore[user.User]{
		localStore: map[string]*user.User{
			"1": {
				ID:                "1",
				ConversationUUIDs: []string{},
			},
		},
		makeDeepCopy: func(usr *user.User) *user.User {
			userCopy := user.User{
				ID:                usr.ID,
				ConversationUUIDs: []string{},
			}

			// deep copy
			userCopy.ConversationUUIDs = append(userCopy.ConversationUUIDs, usr.ConversationUUIDs...)

			return &userCopy
		},
		getID: func(usr user.User) string {
			return usr.ID
		},
	}

	mockAudioClipStore = &mockAssetStore{}
)

func handleStartArtConversation(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     r.Method,
		Body:           buf.String(),
		RequestContext: testUserRequestCtx,
	}

	ctx := startartconversation.HandlerCtx{
		ConversationStore: mockConversationStore,
		UserStore:         mockUserStorage,
		AudioClipStore:    mockAudioClipStore,
		TextGen:           mockTextGen,
		SpeechGen:         mockSpeechGen,
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

func handleContinueConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod: r.Method,
		PathParameters: map[string]string{
			"uuid": vars["id"],
		},
		Body:           buf.String(),
		RequestContext: testUserRequestCtx,
	}

	ctx := continueconversation.HandlerCtx{
		ConversationStore: mockConversationStore,
		UserStore:         mockUserStorage,
		AudioClipStore:    mockAudioClipStore,
		TextGen:           mockTextGen,
		SpeechGen:         mockSpeechGen,
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

func handleGetConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	event := events.APIGatewayProxyRequest{
		HTTPMethod: r.Method,
		PathParameters: map[string]string{
			"uuid": vars["id"],
		},
		RequestContext: testUserRequestCtx,
	}

	ctx := getconversation.HandlerCtx{
		ConversationStore: mockConversationStore,
		UserStore:         mockUserStorage,
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

func handleListConversations(w http.ResponseWriter, r *http.Request) {
	event := events.APIGatewayProxyRequest{
		HTTPMethod:     r.Method,
		RequestContext: testUserRequestCtx,
	}

	ctx := listconversations.HandlerCtx{
		ConversationStore: mockConversationStore,
		UserStore:         mockUserStorage,
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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	router := mux.NewRouter()

	router.HandleFunc("/api/conversation/create/art", handleStartArtConversation).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/list", handleListConversations).Methods(http.MethodGet)
	router.HandleFunc("/api/conversation/{id}/continue", handleContinueConversation).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/{id}", handleGetConversation).Methods(http.MethodGet)
	router.HandleFunc("/api/app/logs", handleAppLogs).Methods(http.MethodPost)
	router.HandleFunc("/aws/presigned/{id}", fileHandler).Methods(http.MethodGet)

	port := ":8080"
	log.Info().Msgf("Server running on port %s", port)
	log.Fatal().Err(http.ListenAndServe(port, router))
}
