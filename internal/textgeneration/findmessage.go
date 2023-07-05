package textgeneration

func (conversation Conversation) FindLastMessageBy(role string) *Message {
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
