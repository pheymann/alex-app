package testutil

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"talktome.com/internal/conversation"
	"talktome.com/internal/processqueue"
	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
	"talktome.com/internal/user"
)

type MockTextGeneration struct {
	GeneratedMessage string
	Timeout          time.Duration
}

func (mock *MockTextGeneration) GenerateNextMessage(
	messageHistory []textgeneration.BasicMessage,
	logCtx zerolog.Context,
) (*textgeneration.BasicMessage, error) {
	// simulate OpenAI generation time
	time.Sleep(mock.Timeout)

	return &textgeneration.BasicMessage{
		Role: textgeneration.RoleAssistent,
		Text: mock.GeneratedMessage,
	}, nil
}

type MockSpeechGeneration struct {
	TestFile string
	Timeout  time.Duration
}

func (mock *MockSpeechGeneration) GenerateSpeechClip(
	title string,
	text string,
	language shared.Language,
	logCtx zerolog.Context,
) (*os.File, error) {
	file, err := os.Open(mock.TestFile)
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
	time.Sleep(mock.Timeout)

	return copyFile, nil
}

func MockConversationStore(initData map[string]*conversation.Conversation) *MockEntityStore[conversation.Conversation] {
	return &MockEntityStore[conversation.Conversation]{
		LocalStore: initData,
		MakeDeepCopy: func(conv *conversation.Conversation) *conversation.Conversation {
			convCopy := conversation.Conversation{
				ID:        conv.ID,
				Metadata:  conv.Metadata,
				Messages:  []conversation.Message{},
				State:     conv.State,
				CreatedAt: conv.CreatedAt,
			}

			// deep copy
			convCopy.Messages = append(convCopy.Messages, conv.Messages...)

			return &convCopy
		},
		GetID: func(conv conversation.Conversation) string {
			return conv.ID
		},
	}
}

func MockUserStore(initData map[string]*user.User) *MockEntityStore[user.User] {
	return &MockEntityStore[user.User]{
		LocalStore: initData,
		MakeDeepCopy: func(usr *user.User) *user.User {
			userCopy := user.User{
				ID:                usr.ID,
				ConversationUUIDs: []string{},
			}

			// deep copy
			userCopy.ConversationUUIDs = append(userCopy.ConversationUUIDs, usr.ConversationUUIDs...)

			return &userCopy
		},
		GetID: func(usr user.User) string {
			return usr.ID
		},
	}
}

type MockEntityStore[E any] struct {
	LocalStore   map[string]*E
	MakeDeepCopy func(entity *E) *E
	GetID        func(entity E) string
}

func (mock *MockEntityStore[E]) Find(uuid string, logCtx zerolog.Context) (*E, error) {
	if entity, ok := mock.LocalStore[uuid]; ok {
		return mock.MakeDeepCopy(entity), nil
	}
	return nil, nil
}

func (mock *MockEntityStore[E]) FindAll(uuids []string, logCtx zerolog.Context) ([]E, error) {
	var entities []E
	for _, uuid := range uuids {
		if entity, ok := mock.LocalStore[uuid]; ok {
			entities = append(entities, *entity)
		}
	}

	return entities, nil
}

func (mock *MockEntityStore[E]) Save(entity E, logCtx zerolog.Context) error {
	mock.LocalStore[mock.GetID(entity)] = mock.MakeDeepCopy(&entity)
	return nil
}

type MockAssetStore struct {
	ClipKey      string
	PresignedUrl string
}

func (mock *MockAssetStore) Save(file *os.File, logCtx zerolog.Context) (string, error) {
	return mock.ClipKey, nil
}

func (mock *MockAssetStore) GenerateTemporaryAccessURL(audioClipUUID string, logCtx zerolog.Context) (string, error) {
	return mock.PresignedUrl, nil
}

type MockProcessQueue struct {
	Queue chan processqueue.Task
}

func (mock *MockProcessQueue) Enqueue(task processqueue.Task, logCtx zerolog.Context) error {
	mock.Queue <- task
	return nil
}

type MockIDGenerator struct {
	GeneratedID string
	UseMetadata bool
}

func (mock *MockIDGenerator) GenerateID(metadata map[string]string) string {
	if mock.UseMetadata {
		return fmt.Sprintf("%s::%+v", mock.GeneratedID, metadata)
	}
	return mock.GeneratedID
}
