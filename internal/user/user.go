package user

type User struct {
	ID                string   `json:"id" yaml:"id" dynamodbav:"id"`
	ConversationUUIDs []string `json:"conversationUuids" yaml:"conversationUuids" dynamodbav:"conversation_uuids"`
}

func (user User) HasConversation(convUUID string) bool {
	for _, uuid := range user.ConversationUUIDs {
		if uuid == convUUID {
			return true
		}
	}

	return false
}
