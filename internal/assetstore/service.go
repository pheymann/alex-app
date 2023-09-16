package assetstore

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type AssetStore interface {
	Save(file *os.File, logCtx zerolog.Context) (string, error)
	GenerateTemporaryAccessURL(key string, logCtx zerolog.Context) (string, *time.Time, error)
}
