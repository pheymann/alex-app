package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/art"
	"talktome.com/internal/cmd/talktome"
	"talktome.com/internal/textgeneration"
)

type mockTextGeneration struct{}

func (generator *mockTextGeneration) ContinueConversation(conversation *textgeneration.Conversation) error {
	conversation.Messages = append(conversation.Messages, textgeneration.Message{
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

	return file, nil
}

type mockStorageService struct{}

func (ctx *mockStorageService) FindArtConversation(uuid string) (*art.ArtConversation, error) {
	return nil, nil
}

func (ctx *mockStorageService) StoreArtConversation(conversation art.ArtConversation) error {
	return nil
}

func (ctx *mockStorageService) StoreClip(clipFile *os.File) error {
	return nil
}

var (
	ttm = talktome.NewTalkToMe(&mockTextGeneration{}, &mockSpeechGeneration{}, &mockStorageService{})
)

func artHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var artPiece art.ArtPiece
		err := json.NewDecoder(r.Body).Decode(&artPiece)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf(">> POST /api/art: %+v\n", artPiece)

		conversation, _, err := ttm.TalkToMeArt(artPiece, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conversation.ConversationStartClipUUID = "prompt.mp3"

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(conversation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Print(">> GET /api/assets: response = prompt.mp3\n")

		http.ServeFile(w, r, "assets/prompt.mp3")
	}
}

func main() {

	http.HandleFunc("/api/art", artHandler)
	http.HandleFunc("/api/assets/", fileHandler)

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
