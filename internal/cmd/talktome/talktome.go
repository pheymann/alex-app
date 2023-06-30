package talktome

import (
	"fmt"

	"github.com/resemble-ai/resemble-go/v2/response"
	"talktome.com/internal/art"
	"talktome.com/internal/speechgeneration"
	"talktome.com/internal/textgeneration"
)

type TalkToMe struct {
	textGen    *textgeneration.TextGenerator
	speechGen  *speechgeneration.SpeechGenerator
	artStorage *art.StorageCtx
}

func NewTalkToMe(textGen *textgeneration.TextGenerator, speechGen *speechgeneration.SpeechGenerator, storage *art.StorageCtx) TalkToMe {
	return TalkToMe{
		textGen:    textGen,
		speechGen:  speechGen,
		artStorage: storage,
	}
}

var (
	emptyPresentation = art.ArtPresentation{}
)

func (talktome TalkToMe) GetOrCreatePresentation(piece art.ArtPiece) (art.ArtPresentation, error) {
	presentationID := fmt.Sprintf("%s::%s", piece.ArtistName, piece.Name)
	presentation, err := talktome.artStorage.FindArtPresentation(presentationID)
	if err != nil {
		return emptyPresentation, err
	}

	// no presentation found
	if presentation == nil {
		return talktome.generatePresentation(piece)
	}

	return *presentation, nil
}

func (talktome TalkToMe) generatePresentation(piece art.ArtPiece) (art.ArtPresentation, error) {
	presentation, err := talktome.generateTextContent(piece)
	if err != nil {
		return emptyPresentation, err
	}

	clip, err := talktome.generateSpeechClip(presentation.Description)
	if err != nil {
		return emptyPresentation, err
	}

	presentation.DescriptionClipURL = clip.Item.AudioSrc
	if err := talktome.artStorage.StoreArtPresentation(presentation); err != nil {
		return emptyPresentation, err
	}

	for i, task := range presentation.Tasks {
		clip, err := talktome.generateSpeechClip(task.Task)
		if err != nil {
			// better to continue generating clips and having some content instead of
			// failing all attempts here and leave the user with nothing in the worst case
			fmt.Printf("[WARN] failed to generate clip for task %s: %s", task.Task, err)
		}

		presentation.Tasks[i].TaskClipURL = clip.Item.AudioSrc
		if err := talktome.artStorage.StoreArtPresentation(presentation); err != nil {
			// same reasoning here
			fmt.Printf("[WARN] failed to store clip for task %s: %s", task.Task, err)
		}
	}

	return presentation, nil
}

func (talktome TalkToMe) generateTextContent(piece art.ArtPiece) (art.ArtPresentation, error) {
	fmt.Printf("[DEBUG] Generate description for %s's \"%s\"\n", piece.ArtistName, piece.Name)

	description, err := talktome.textGen.GenerateArtDescription(piece.ArtistName, piece.Name)
	if err != nil {
		return emptyPresentation, err
	}

	fmt.Printf("[DEBUG] Generate tasks for %s's \"%s\"\n", piece.ArtistName, piece.Name)

	taskTexts, err := talktome.textGen.GenerateTasks(piece.ArtistName, piece.Name)
	if err != nil {
		return emptyPresentation, err
	}

	var tasks = []art.ArtPresentationTask{}
	for _, text := range taskTexts {
		tasks = append(tasks, art.ArtPresentationTask{
			Task: text,
		})
	}

	return art.ArtPresentation{
		ID:          fmt.Sprintf("%s::%s", piece.ArtistName, piece.Name),
		Piece:       piece,
		Description: description,
		Tasks:       tasks,
	}, nil
}

func (talktome TalkToMe) generateSpeechClip(text string) (response.Clip, error) {
	fmt.Printf("[DEBUG] Generate speech clip for %s\n", text)
	return talktome.speechGen.GenerateSpeechClip(text)
}
