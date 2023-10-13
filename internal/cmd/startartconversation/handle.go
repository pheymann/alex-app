package startartconversation

import (
	"fmt"

	"talktome.com/internal/conversation"
	"talktome.com/internal/shared"
	"talktome.com/internal/textgeneration"
)

const (
	MaxArtContextLength = 150
)

func Handle(ctx conversation.Context, artContext string) (*conversation.Conversation, error) {
	if artContext == "" {
		return nil, &shared.UserInputError{Message: "artist context cannot be empty"}
	} else if len(artContext) > MaxArtContextLength {
		shared.GetLogger(ctx.LogCtx).Warn().Msgf("artist context too long: %d", len(artContext))
		return nil, &shared.UserInputError{Message: "artist context too long"}
	}

	shared.GetLogger(ctx.LogCtx).Debug().Msgf("start art conversation for '%s'", artContext)

	metadata := map[string]string{
		"artContext": artContext,
	}
	conv := conversation.NewConversation(metadata)
	ctx.LogCtx = ctx.LogCtx.Str("conversation_uuid", conv.ID)

	conv.Messages = []conversation.Message{
		{
			Role: textgeneration.RoleSystem,
			Text: `You are an art museum curator and show and explain art pieces to a visitor.
							You have an engaging, friendly, and professional communication style. You talk to a single person and
							you already discussed a couple of paintings. So that is not the beginning of this conversation.
							Finally, you address the visitor with the word "you". Start your first response with one of the following
							phrases: "Here, we stand in front of", "This is", "We are looking at", "This is a painting by".
							Don't use more than 200 words in all your responses. If you don't know the answer say so. If I ask for
							something not related to art or this art piece say "I don't know" or "I don't understand".`,
		},
	}

	ctx.ConversationUUID = conv.ID

	return ctx.StartConversation(
		conv,
		fmt.Sprintf(`The art piece we are discussion is "%s".`, artContext)+
			`Introduce it to me, give some basic information like the creation date, and continue to explain
			its meaning, what style it is, and how it fits into its time.`,
	)
}
