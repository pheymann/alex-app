package conversation

import "time"

type ConversationRef struct {
	ID        string            `json:"id"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"createdAt"`
}
