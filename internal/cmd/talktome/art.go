package talktome

import (
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
	"talktome.com/internal/art"
)

func (talktome TalkToMe) TalkToMeArt(piece art.ArtPiece, prompt *string) (*art.ArtConversation, []byte, error) {
	fmt.Printf("[DEBUG] load conversation for %s\n", piece.String())
	lookedUpConversation, err := talktome.artStorage.FindArtConversation(art.CreateArtConversationID(piece))
	if err != nil {
		return nil, nil, err
	}

	var conversation art.ArtConversation

	// no conversation found
	if lookedUpConversation == nil {
		fmt.Printf("[DEBUG] start new conversation for %s\n", piece.String())
		newConversation, err := art.StartArtConversation(talktome.textGen, piece)
		if err != nil {
			return nil, nil, err
		}

		fmt.Printf("[DEBUG] store conversation %s\n", conversation.ID)
		if err := talktome.artStorage.StoreArtConversation(conversation); err != nil {
			return nil, nil, err
		}

		conversation = *newConversation
	} else {
		conversation = *lookedUpConversation
	}

	fmt.Printf("[DEBUG] conversation about %s continues with ID %s\n", piece.String(), conversation.ID)

	// checking the existence of the clip UUIDs to cover the case where generating text worked but something
	// broke during clip creation last time we tried
	if conversation.ConversationStartClipUUID == "" {
		clipFile, err := talktome.generateClip(conversation)
		if err != nil {
			return nil, nil, err
		}

		defer clipFile.Close()

		fmt.Printf("[DEBUG] store clip for conversation %s\n", conversation.ID)
		if err := talktome.artStorage.StoreClip(clipFile); err != nil {
			return nil, nil, err
		}

		fmt.Printf("[DEBUG] store conversation %s\n", conversation.ID)
		if err := talktome.artStorage.StoreArtConversation(conversation); err != nil {
			return nil, nil, err
		}

		conversation.ConversationStartClipUUID = clipFile.Name()
	}

	var promptClip []byte

	if prompt != nil {
		fmt.Printf("[DEBUG] continue with user prompt for conversation %s\n", conversation.ID)

		conversation.ConversationStart.AddPrompt(*prompt)

		err := talktome.textGen.ContinueConversation(&conversation.ConversationStart)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to continue conversation %s: %w", conversation.ID, err)
		}

		fmt.Printf("[DEBUG] generate clip for continuation for conversation %s\n", conversation.ID)
		clipFile, err := talktome.generateClip(conversation)
		if err != nil {
			return nil, nil, err
		}

		defer clipFile.Close()

		fileInfo, err := clipFile.Stat()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get clip file statistics:%w", err)
		}

		promptClip = make([]byte, fileInfo.Size())

		_, err = io.ReadFull(clipFile, promptClip)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil, nil, fmt.Errorf("failed to read clip bytes:%w", err)
		}
	}

	// don't show initialization prompts
	conversation.ConversationStart.Messages = conversation.ConversationStart.Messages[3:]

	return &conversation, promptClip, nil
}

func (talktome TalkToMe) generateClip(conversation art.ArtConversation) (*os.File, error) {
	fmt.Printf("[DEBUG] generate clip for conversation %s\n", conversation.ID)
	message := conversation.ConversationStart.FindLastMessageBy(openai.ChatMessageRoleAssistant)

	if message != nil {
		clipFile, err := talktome.speechGen.GenerateSpeechClip(conversation.ID, message.Text)
		if err != nil {
			return nil, err
		}

		return clipFile, nil
	} else {
		return nil, fmt.Errorf("missing assistent message")
	}
}
