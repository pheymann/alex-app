package conversation

type ConversationRef struct {
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata"`
}
