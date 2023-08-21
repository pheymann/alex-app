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

func (generator *mockTextGeneration) GenerateNextMessage(conv *conversation.Conversation) error {
	conv.Messages = append(conv.Messages, conversation.Message{
		Role: openai.ChatMessageRoleAssistant,
		Text: "A long explanation ....",
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod:     r.Method,
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
	router.HandleFunc("/api/conversation/continue", handleContinueConversation).Methods(http.MethodPost)
	router.HandleFunc("/api/conversation/list", handleListConversations).Methods(http.MethodGet)
	router.HandleFunc("/api/conversation/{id}", handleGetConversation).Methods(http.MethodGet)
	router.HandleFunc("/aws/presigned/{id}", fileHandler).Methods(http.MethodGet)

	port := ":8080"
	log.Info().Msgf("Server running on port %s", port)
	log.Fatal().Err(http.ListenAndServe(port, router))
}
