package conversation

func (conversation Conversation) FindLastMessageBy(role Role) *Message {
	lastIndex := len(conversation.Messages) - 1

	for {
		if conversation.Messages[lastIndex].Role == role {
			return &conversation.Messages[lastIndex]
		}

		lastIndex--

		if lastIndex < 0 {
			return nil
		}
	}
}
