package user

type User struct {
	ID                string   `json:"id" dynamodbav:"id"`
	ConversationUUIDs []string `json:"conversation_uuids" dynamodbav:"conversation_uuids"`
}
