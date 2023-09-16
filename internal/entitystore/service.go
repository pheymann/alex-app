package entitystore

import "github.com/rs/zerolog"

type EntityStore[E any] interface {
	Find(uuid string, logCtx zerolog.Context) (*E, error)
	FindAll(uuids []string, logCtx zerolog.Context) ([]E, error)
	Save(entity E, logCtx zerolog.Context) error
}
