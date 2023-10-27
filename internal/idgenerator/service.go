package idgenerator

type IDGenerator interface {
	GenerateID(metadata map[string]string) string
}
