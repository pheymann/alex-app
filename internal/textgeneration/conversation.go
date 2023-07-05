package textgeneration

type Conversation struct {
	Messages []Message
}

type Message struct {
	Role string
	Text string
}
