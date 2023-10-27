package idgenerator

type randomIDGenerator struct{}

func NewRandomIDGenerator() IDGenerator {
	return &randomIDGenerator{}
}
