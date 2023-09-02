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
	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/cmd/continueconversation"
	"talktome.com/internal/cmd/getconversation"
	"talktome.com/internal/cmd/listconversations"
	"talktome.com/internal/cmd/startartconversation"
	"talktome.com/internal/conversation"
	"talktome.com/internal/talktome"
	"talktome.com/internal/user"
)

type mockTextGeneration struct{}

const (
	longText = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.

	Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.

	Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.

	Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis.

	At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, At accusam aliquyam diam diam dolore dolores duo eirmod eos erat, et nonumy sed tempor et et invidunt justo labore Stet clita ea et gubergren, kasd magna no rebum. sanctus sea sed takimata ut vero voluptua. est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur`
)

func (generator *mockTextGeneration) GenerateNextMessage(conv *conversation.Conversation) error {
	conv.Messages = append(conv.Messages, conversation.Message{
		Role: openai.ChatMessageRoleAssistant,
		Text: longText,
	})

	// simulate OpenAI generation time
	time.Sleep(5 * time.Second)

	return nil
}

type mockSpeechGeneration struct{}

func (generator *mockSpeechGeneration) GenerateSpeechClip(title string, text string) (*os.File, error) {
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

type mockConversationStorageService struct {
	storage map[string]*conversation.Conversation
}

func (ctx *mockConversationStorageService) FindConversation(uuid string) (*conversation.Conversation, error) {
	if conv, ok := ctx.storage[uuid]; ok {
		convCopy := conversation.Conversation{
			ID:       conv.ID,
			Metadata: conv.Metadata,
			Messages: []conversation.Message{},
		}

		// deep copy
		convCopy.Messages = append(convCopy.Messages, conv.Messages...)

		return &convCopy, nil
	}
	return nil, nil
}

func (ctx *mockConversationStorageService) FindAllConversations(uuids []string) ([]conversation.Conversation, error) {
	var conversations []conversation.Conversation
	for _, uuid := range uuids {
		if conversation, ok := ctx.storage[uuid]; ok {
			conversations = append(conversations, *conversation)
		}
	}

	return conversations, nil
}

func (ctx *mockConversationStorageService) StoreConversation(conv conversation.Conversation) error {
	convCopy := conversation.Conversation{
		ID:       conv.ID,
		Metadata: conv.Metadata,
		Messages: []conversation.Message{},
	}

	// deep copy
	convCopy.Messages = append(convCopy.Messages, conv.Messages...)

	ctx.storage[conv.ID] = &convCopy
	return nil
}

func (ctx *mockConversationStorageService) StoreClip(clipFile *os.File) (string, error) {
	return "prompt.mp3", nil
}

func (ctx *mockConversationStorageService) GenerateClipAccess(audioClipUUID string) (string, error) {
	return "/aws/presigned/prompt.mp3", nil
}

type mockUserStorageService struct {
	storage map[string]*user.User
}

func (ctx *mockUserStorageService) FindUser(uuid string) (*user.User, error) {
	if user, ok := ctx.storage[uuid]; ok {
		return user, nil
	}
	return nil, nil
}

func (ctx *mockUserStorageService) StoreUser(user user.User) error {
	ctx.storage[user.ID] = &user
	return nil
}

var (
	mockTextGen     = &mockTextGeneration{}
	mockSpeechGen   = &mockSpeechGeneration{}
	mockConvStorage = &mockConversationStorageService{storage: map[string]*conversation.Conversation{}}
	mockUserStorage = &mockUserStorageService{storage: map[string]*user.User{
		"1": {
			ID: "1",
			ConversationUUIDs: []string{
				"1",
			},
		},
	}}

	mockCtx = talktome.NewContext(mockTextGen, mockSpeechGen, mockConvStorage, mockUserStorage)
)

var (
	testRequestContext = events.APIGatewayProxyRequestContext{
		Authorizer: map[string]interface{}{
			"jwt": map[string]interface{}{
				"claims": map[string]interface{}{
					"cognito:username": "1",
				},
			},
		},
	}
)

func handleCreateArtConversation(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     r.Method,
		Body:           buf.String(),
		RequestContext: testRequestContext,
	}

	response, err := startartconversation.HandlerCtx{Ctx: mockCtx}.AWSHandler(context.TODO(), event)
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
		RequestContext: testRequestContext,
	}

	response, err := continueconversation.HandlerCtx{Ctx: mockCtx}.AWSHandler(context.TODO(), event)
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
		RequestContext: testRequestContext,
	}

	response, err := getconversation.HandlerCtx{UserStorage: mockUserStorage, ConvStorage: mockConvStorage}.AWSHandler(context.TODO(), event)
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
		RequestContext: testRequestContext,
	}

	response, err := listconversations.HandlerCtx{UserStorage: mockUserStorage, ConvStorage: mockConvStorage}.AWSHandler(context.TODO(), event)
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

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "assets/prompt.mp3")
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	router := mux.NewRouter()

	router.HandleFunc("/api/conversation/create/art", handleCreateArtConversation).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/list", handleListConversations).Methods(http.MethodGet)
	router.HandleFunc("/api/conversation/{id}/continue", handleContinueConversation).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/{id}", handleGetConversation).Methods(http.MethodGet)
	router.HandleFunc("/aws/presigned/{id}", fileHandler).Methods(http.MethodGet)

	port := ":8080"
	log.Info().Msgf("Server running on port %s", port)
	log.Fatal().Err(http.ListenAndServe(port, router))
}
