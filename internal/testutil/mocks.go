package testutil

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"talktome.com/internal/textgeneration"
)

type MockTextGeneration struct {
	GeneratedMessage string
}

func (mock *MockTextGeneration) GenerateNextMessage(
	messageHistory []textgeneration.BasicMessage,
	logCtx zerolog.Context,
) (*textgeneration.BasicMessage, error) {
	// simulate OpenAI generation time
	time.Sleep(5 * time.Second)

	return &textgeneration.BasicMessage{
		Role: textgeneration.RoleAssistent,
		Text: mock.GeneratedMessage,
	}, nil
}

type MockSpeechGeneration struct {
	TestFile string
}

func (mock *MockSpeechGeneration) GenerateSpeechClip(
	title string,
	text string,
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
	time.Sleep(3 * time.Second)

	return copyFile, nil
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

type MockAssetStore struct{}

func (mock *MockAssetStore) Save(file *os.File, logCtx zerolog.Context) (string, error) {
	return "prompt.mp3", nil
}

func (mock *MockAssetStore) GenerateTemporaryAccessURL(audioClipUUID string, logCtx zerolog.Context) (string, *time.Time, error) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return "", nil, err
	}

	urlValidFor := 1 * time.Minute
	expirationDate := time.Now().In(location).Add(urlValidFor)

	return "/aws/presigned/prompt.mp3", &expirationDate, nil
}
