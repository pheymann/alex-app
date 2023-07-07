package art

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/textgeneration"
)

func StartArtConversation(textGen textgeneration.TextGenerator, piece ArtPiece) (*ArtConversation, error) {
	start := textgeneration.Conversation{
		Messages: []textgeneration.Message{
			{
				Role: openai.ChatMessageRoleSystem,
				Text: fmt.Sprintf(`The art piece we are discussion is "%s" from %s`, piece.ArtistName, piece.Name),
			},
			{
				Role: openai.ChatMessageRoleSystem,
				Text: `You are an art museum curator and show and explain art pieces to a visitor. You have an engaging, friendly, and professional communication style. You talk to a single person and you already discussed a couple of paintings already. So that is not the beginning of this conversation. Finally, you address the visitor with the word "you".`,
			},
			{
				Role: openai.ChatMessageRoleUser,
				Text: `We are standing in front of the art piece. Introduce it to me, give some basic information like the creation date, and continue to explain its meaning, what style it is, and how it fits into its time. Don't use more than 200 words.`,
			},
		},
	}

	err := textGen.ContinueConversation(&start)
	if err != nil {
		return nil, fmt.Errorf("failed to start conversation for %s's \"%s\": %w", piece.ArtistName, piece.Name, err)
	}

	conversation := ArtConversation{
		ID:                CreateArtConversationID(piece),
		Piece:             piece,
		ConversationStart: start,
	}

	return &conversation, nil
}
