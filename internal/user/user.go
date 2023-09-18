package user

type User struct {
	ID                string   `json:"id" yaml:"id" dynamodbav:"id"`
	ConversationUUIDs []string `json:"conversationUuids" yaml:"conversationUuids" dynamodbav:"conversation_uuids"`
}
