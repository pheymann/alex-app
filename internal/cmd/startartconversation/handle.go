package startartconversation

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/conversation"
	"talktome.com/internal/talktome"
)

func Handle(userUUID string, artContext string, ctx talktome.Context) (*conversation.Conversation, error) {
	log.Info().Str("user_uuid", userUUID).Msgf("start art conversation for '%s'", artContext)

	if artContext == "" {
		return nil, fmt.Errorf("artist context cannot be empty")
	}

	metadata := map[string]string{
		"artContext": artContext,
	}
	conv := conversation.NewConversation(metadata)
	conv.Messages = []conversation.Message{
		{
			Role:        openai.ChatMessageRoleSystem,
			Text:        fmt.Sprintf(`The art piece we are discussion is "%s"`, artContext),
			CanHaveClip: false,
		},
		{
			Role:        openai.ChatMessageRoleSystem,
			Text:        `You are an art museum curator and show and explain art pieces to a visitor. You have an engaging, friendly, and professional communication style. You talk to a single person and you already discussed a couple of paintings already. So that is not the beginning of this conversation. Finally, you address the visitor with the word "you".`,
			CanHaveClip: false,
		},
	}

	return ctx.StartConversation(
		userUUID,
		conv,
		`We are standing in front of the art piece. Introduce it to me, give some basic information like the creation date, and continue to explain its meaning, what style it is, and how it fits into its time. Don't use more than 200 words.`,
	)
}
