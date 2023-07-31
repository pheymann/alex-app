package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/cmd/talktomeartcreate"
	"talktome.com/internal/cmd/talktomecontinue"
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
	return nil
}

type mockSpeechGeneration struct{}

func (generator *mockSpeechGeneration) GenerateSpeechClip(title string, text string) (*os.File, error) {
	file, err := os.Open("assets/prompt.mp3")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	copyFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(copyFile, file); err != nil {
		return nil, err
	}

	return copyFile, nil
}

type mockConversationStorageService struct {
	storage map[string]*conversation.Conversation
}

func (ctx *mockConversationStorageService) FindConversation(uuid string) (*conversation.Conversation, error) {
	if conversation, ok := ctx.storage[uuid]; ok {
		return conversation, nil
	}
	return nil, nil
}

func (ctx *mockConversationStorageService) StoreConversation(conv conversation.Conversation) error {
	ctx.storage[conv.ID] = &conv
	return nil
}

func (ctx *mockConversationStorageService) StoreClip(clipFile *os.File) error {
	return nil
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
	mockConvStorage = &mockConversationStorageService{storage: make(map[string]*conversation.Conversation)}
	mockUserStorage = &mockUserStorageService{storage: map[string]*user.User{
		"1": {
			ID: "1",
		},
	}}

	mockCtx = talktome.NewContext(mockTextGen, mockSpeechGen, mockConvStorage, mockUserStorage)
)

func handleCreateArtConversation(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	event := events.APIGatewayProxyRequest{
		HTTPMethod: r.Method,
		Body:       buf.String(),
	}

	response, err := talktomeartcreate.HandlerCtx{Ctx: mockCtx}.AWSHandler(context.TODO(), event)
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
		HTTPMethod: r.Method,
		Body:       buf.String(),
	}

	response, err := talktomecontinue.HandlerCtx{Ctx: mockCtx}.AWSHandler(context.TODO(), event)
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
		log.Info().Msgf(">> GET %s", r.URL.Path)

		http.ServeFile(w, r, "assets/prompt.mp3")
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	http.HandleFunc("/api/conversation/create/art", handleCreateArtConversation)
	http.HandleFunc("/api/conversation/continue", handleContinueConversation)
	http.HandleFunc("/api/assets/", fileHandler)

	port := ":8080"
	log.Info().Msgf("Server running on port %s", port)
	log.Fatal().Err(http.ListenAndServe(port, nil))
}