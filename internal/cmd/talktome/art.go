package talktome

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/art"
)

func (talktome TalkToMe) TalkToMeArt(piece art.ArtPiece) (*art.ArtConversation, error) {
	fmt.Printf("[DEBUG] load conversation for %s\n", piece.String())
	lookedUpConversation, err := talktome.artStorage.FindArtConversation(art.CreateArtConversationID(piece))
	if err != nil {
		return nil, err
	}

	var conversation art.ArtConversation

	// no conversation found
	if lookedUpConversation == nil {
		fmt.Printf("[DEBUG] start new conversation for %s\n", piece.String())
		newConversation, err := art.StartArtConversation(talktome.textGen, piece)
		if err != nil {
			return nil, err
		}

		conversation = *newConversation
	} else {
		conversation = *lookedUpConversation
	}

	fmt.Printf("[DEBUG] conversation about %s continues with ID %s\n", piece.String(), conversation.ID)

	fmt.Printf("[DEBUG] store conversation %s\n", conversation.ID)
	if err := talktome.artStorage.StoreArtConversation(conversation); err != nil {
		return nil, err
	}

	// checking the existence of the clip UUIDs to cover the case where generating text worked but something
	// broke during clip creation last time we tried
	if conversation.ConversationStartClipUUID == "" {
		fmt.Printf("[DEBUG] generate clip for conversation %s\n", conversation.ID)
		message := conversation.ConversationStart.FindLastMessageBy(openai.ChatMessageRoleAssistant)

		if message != nil {
			clipFile, err := talktome.speechGen.GenerateSpeechClip(conversation.ID, message.Text)
			if err != nil {
				return nil, err
			}

			defer clipFile.Close()

			fmt.Printf("[DEBUG] store clip for conversation %s\n", conversation.ID)
			if err := talktome.artStorage.StoreClip(clipFile); err != nil {
				return nil, err
			}

			conversation.ConversationStartClipUUID = clipFile.Name()

			fmt.Printf("[DEBUG] store conversation %s\n", conversation.ID)
			if err := talktome.artStorage.StoreArtConversation(conversation); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("missing assistent message")
		}
	}

	return &conversation, nil
}
