package shared

import (
	"github.com/rs/zerolog"
)

func GetLogger(logCtx zerolog.Context) *zerolog.Logger {
	logger := logCtx.Logger()

	return &logger
}
