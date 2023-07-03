package talktome

import (
	"fmt"

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
	lookupPresentation, err := talktome.artStorage.FindArtPresentation(art.CreatePresentationID(piece))
	if err != nil {
		return emptyPresentation, err
	}

	var presentation art.ArtPresentation

	// no presentation found
	if lookupPresentation == nil {
		presentation, err = talktome.generateTextContent(piece)
		if err != nil {
			return emptyPresentation, err
		}
	} else {
		presentation = *lookupPresentation
	}

	// checking the existence of the clip UUIDs to cover the case where generating text worked but something
	// broke during clip creation last time we tried
	if presentation.DescriptionClipUUID == "" {
		err := talktome.generateAndStoreClip(presentation.ID, presentation.Description, &presentation, func(presentation *art.ArtPresentation, uuid string) {
			presentation.DescriptionClipUUID = uuid
		})
		if err != nil {
			return emptyPresentation, err
		}
	}

	// for index, task := range presentation.Tasks {
	// 	if task.TaskClipUUID == "" {
	// 		talktome.generateAndStoreClip(presentation.ID, presentation.Description, &presentation, func(presentation *art.ArtPresentation, uuid string) {
	// 			presentation.Tasks[index].TaskClipUUID = uuid
	// 		})
	// 	}
	// }

	return presentation, nil
}

func (talktome TalkToMe) generateAndStoreClip(
	title,
	text string,
	presentation *art.ArtPresentation,
	updateUUID func(*art.ArtPresentation, string),
) error {
	fmt.Printf("[DEBUG] Generate clip audio file for %s\n", title)
	clipFile, err := talktome.speechGen.GenerateSpeechClip(title, text)
	if err != nil {
		return err
	}

	defer clipFile.Close()

	fmt.Printf("[DEBUG] Store clip audio file for %s\n", title)
	if err := talktome.artStorage.StoreClip(clipFile); err != nil {
		return err
	}

	fmt.Printf("[DEBUG] Store presentation in database for %s\n", title)
	updateUUID(presentation, clipFile.Name())
	if err := talktome.artStorage.StoreArtPresentation(*presentation); err != nil {
		return err
	}

	return nil
}

func (talktome TalkToMe) generateTextContent(piece art.ArtPiece) (art.ArtPresentation, error) {
	fmt.Printf("[DEBUG] Generate description for %s's \"%s\"\n", piece.ArtistName, piece.Name)

	description, err := talktome.textGen.GenerateArtDescription(piece.ArtistName, piece.Name)
	if err != nil {
		return emptyPresentation, err
	}

	fmt.Printf("[DEBUG] Generate tasks for %s's \"%s\"\n", piece.ArtistName, piece.Name)

	// taskTexts, err := talktome.textGen.GenerateTasks(piece.ArtistName, piece.Name)
	// if err != nil {
	// 	return emptyPresentation, err
	// }

	var tasks = []art.ArtPresentationTask{}
	// for _, text := range taskTexts {
	// 	tasks = append(tasks, art.ArtPresentationTask{
	// 		Task: text,
	// 	})
	// }

	return art.ArtPresentation{
		ID:          fmt.Sprintf("%s::%s", piece.ArtistName, piece.Name),
		Piece:       piece,
		Description: description,
		Tasks:       tasks,
	}, nil
}
