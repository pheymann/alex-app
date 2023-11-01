package startartconversation

import (
	"fmt"

	"talktome.com/internal/conversation"
	"talktome.com/internal/processqueue"
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
	conv := conversation.NewConversation(metadata, ctx.IDGenerator)
	ctx.LogCtx = ctx.LogCtx.Str("conversation_uuid", conv.ID)

	conv.State = conversation.StateGenerating

	systemPrompt, err := systemPrompt(ctx.Language)
	if err != nil {
		return nil, err
	}
	conv.Messages = []conversation.Message{
		{
			Role: textgeneration.RoleSystem,
			Text: systemPrompt,
		},
	}

	ctx.ConversationUUID = conv.ID

	user, err := ctx.UserStore.Find(ctx.UserUUID, ctx.LogCtx)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, &shared.NotFoundError{Message: fmt.Sprintf("user with UUID %s does not exist", ctx.UserUUID)}
	}

	if err := ctx.ConversationStore.Save(conv, ctx.LogCtx); err != nil {
		return nil, err
	}

	user.ConversationUUIDs = append(user.ConversationUUIDs, conv.ID)
	if err := ctx.UserStore.Save(*user, ctx.LogCtx); err != nil {
		return nil, err
	}

	userPromp, err := initialUserPrompt(artContext, ctx.Language)
	if err != nil {
		return nil, err
	}

	task := processqueue.Task{
		ConversationUUID: conv.ID,
		UserUUID:         ctx.UserUUID,
		Language:         ctx.Language,
		Message:          userPromp,
	}

	if err := ctx.ProcessQueue.Enqueue(task, ctx.LogCtx); err != nil {
		return nil, err
	}

	// system messages are not relevant to the client
	conv.Messages = []conversation.Message{}

	return &conv, nil
}

const (
	englishSystemPrompt = `You are an art museum curator and explain art pieces, artists, and art techniques to a visitor.
	You have an engaging, friendly, and professional communication style. You talk to a single person and
		you already discussed a couple of topics. So that is not the beginning of this conversation.
		Finally, you address the visitor with the word "you". Your answers should get immediately to the point with out any introduction.
		Don't use more than 200 words in all your responses. If you don't know the answer say "I don't know". If I ask for
		something not related to art, artists, or art techniques say "I don't know" or "I don't understand".`

	germanSystemPrompt = `Du bist ein Museums Kurator und erklärst Kunstwerke, Künstler und künstlerische Techniken einem Besucher.
	Du hast einen engagierten, freundlichen und professionellen Kommunikationsstil. Du sprichst mit einer einzelnen Person und
		du hast bereits ein paar Themen besprochen. Das ist also nicht der Anfang dieses Gesprächs.
		Zuletzt sprichst du den Besucher mit dem Wort "du" an. Deine Antworten sollten sofort auf den Punkt kommen, ohne jegliche Einleitung.
		Verwende nicht mehr als 200 Wörter in all deinen Antworten. Wenn du die Antwort nicht weißt, sag "Ich weiß es nicht". Wenn ich nach
		etwas frage, das nichts mit Kunst, Künstlern, oder Kunsttechniken zu tun hat, sag "Ich weiß es nicht" oder "Ich verstehe nicht".`
)

func systemPrompt(langauge shared.Language) (string, error) {
	switch langauge {
	case shared.LanguageEnglish:
		return englishSystemPrompt, nil

	case shared.LanguageGerman:
		return germanSystemPrompt, nil

	default:
		return "", &shared.UserInputError{Message: "unsupported language"}
	}
}

func initialUserPrompt(artContext string, language shared.Language) (string, error) {
	switch language {
	case shared.LanguageEnglish:
		return englishUserPrompt(artContext), nil

	case shared.LanguageGerman:
		return germanUserPrompt(artContext), nil

	default:
		return "", &shared.UserInputError{Message: "unsupported language"}
	}
}

func englishUserPrompt(artContext string) string {
	return fmt.Sprintf(`We discuss "%s".`, artContext) + `Give some basic information.`
}

func germanUserPrompt(artContext string) string {
	return fmt.Sprintf(`Wir besprechen "%s".`, artContext) + `Gib einige grundlegende Informationen.`
}
